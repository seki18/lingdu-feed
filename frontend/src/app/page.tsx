'use client';

import { useEffect, useRef, useState } from "react";
import PostCard from "@/components/layout/PostCard";
import { apiFetch, trackState, consumeDirtyPostIds } from "@/lib/api";
import { getUser } from "@/lib/auth";
import { useToast } from "@/components/ui/ToastContext";
import Loading from "@/components/ui/Loading";
import { PostSummary } from "@/types/post";
import { User } from "@/types/user";

// ── Module-level feed cache (survives component remounts) ──
const CACHE_TTL = 30000; // 30 seconds — long enough for back navigation, short enough to not serve stale data
let feedCache: {
  tab: "recommend" | "following";
  posts: PostSummary[];
  deliveredIds: number[];
  hasMore: boolean;
  timestamp: number;
} | null = null;

export default function HomePage() {
  const [tab, setTab] = useState<"recommend" | "following">("recommend");
  const [posts, setPosts] = useState<PostSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [deliveredIds, setDeliveredIds] = useState<number[]>([]);
  const deliveredIdsRef = useRef<number[]>([]);
  const postsRef = useRef<PostSummary[]>([]);
  const initialLoaded = useRef(false);
  const tabLoaded = useRef(false); // track per-tab state
  const [creating, setCreating] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [user, setUser] = useState<User | null>(null);
  const { addToast } = useToast();
  const [showCreateModal, setShowCreateModal] = useState(false);

  // Generic fetch for both recommend and following feeds
  async function fetchFeed(
    endpoint: string,
    requestType: "initial" | "subsequent" = "initial",
    append = false
  ): Promise<number> {
    if (append) setLoadingMore(true);
    else { setLoading(true); setHasMore(true); }

    const currentIdsParam = append && deliveredIdsRef.current.length > 0
      ? deliveredIdsRef.current.join(",")
      : "";

    try {
      const response = await apiFetch(
        `${endpoint}?request_type=${requestType}${currentIdsParam ? `&current_ids=${encodeURIComponent(currentIdsParam)}` : ""}`
      );
      if (response.code !== 200) {
        addToast(response.message || "Failed to load posts.", { type: "error", title: "Load failed" });
        return 0;
      }
      const rawPosts: PostSummary[] = response.data ?? [];
      let newCount = 0;
      if (append) {
        const existingIds = new Set(postsRef.current.map(p => p.id));
        const uniquePosts = rawPosts.filter(p => !existingIds.has(p.id));
        if (uniquePosts.length === 0) { setHasMore(false); return 0; }
        setPosts(current => {
          const next = [...current, ...uniquePosts];
          postsRef.current = next;
          return next;
        });
        newCount = uniquePosts.length;
      } else {
        postsRef.current = rawPosts;
        setPosts(rawPosts);
        newCount = rawPosts.length;
      }
      const fetchedIds = rawPosts.map(p => p.id);
      deliveredIdsRef.current = Array.from(new Set([...deliveredIdsRef.current, ...fetchedIds]));
      setDeliveredIds(deliveredIdsRef.current);
      fetchedIds.forEach(id => trackState(id, 1));
      return newCount;
    } catch (error) {
      addToast("Unable to load posts.", { type: "error", title: "Error" });
      console.error(error);
    } finally {
      if (append) setLoadingMore(false); else setLoading(false);
    }
    return 0;
  }

  // Load current tab
  const loadTab = async (t: "recommend" | "following") => {
    deliveredIdsRef.current = [];
    postsRef.current = [];
    setPosts([]);
    setDeliveredIds([]);
    setHasMore(true);
    hasMoreRef.current = true;
    setObserverPaused(true);
    const endpoint = t === "recommend" ? "/api/feed/recommend" : "/api/feed/following";
    await fetchFeed(endpoint, "initial", false);
    // Fill page, waiting for React render between each batch
    for (let tries = 0; tries < 5; tries++) {
      await new Promise<void>(r => requestAnimationFrame(() => r()));
      if (document.documentElement.scrollHeight > window.innerHeight) break;
      const fetched = await fetchFeed(endpoint, "subsequent", true);
      if (fetched === 0) break;
    }
    // Enable observer after a short delay to let React commit
    setTimeout(() => {
      setObserverPaused(false);
    }, 100);
    // Write to module-level cache
    feedCache = {
      tab: t,
      posts: postsRef.current,
      deliveredIds: deliveredIdsRef.current,
      hasMore: hasMoreRef.current,
      timestamp: Date.now(),
    };
  };

  useEffect(() => {
    const init = async () => {
      if (tabLoaded.current) return;
      tabLoaded.current = true;
      const savedUser = getUser();
      if (savedUser) setUser(savedUser);

      // Restore from module-level cache if fresh enough
      if (feedCache && Date.now() - feedCache.timestamp < CACHE_TTL) {
        setObserverPaused(true);
        setTab(feedCache.tab);
        postsRef.current = feedCache.posts;
        setPosts(feedCache.posts);
        deliveredIdsRef.current = feedCache.deliveredIds;
        setDeliveredIds(feedCache.deliveredIds);
        setHasMore(feedCache.hasMore);
        feedCache.posts.forEach(p => trackState(p.id, 1));
        // Enable observer after React commits
        setTimeout(() => {
          setObserverPaused(false);
        }, 100);

        // If posts were interacted with on the detail page, just invalidate cache;
        // do NOT auto-refresh — user can manually click Refresh to get fresh data
        const ids = consumeDirtyPostIds();
        if (ids.length > 0) {
          feedCache = null;
        }
        return;
      }

      // Clear stale cache before fresh load
      feedCache = null;
      await loadTab("recommend");
    };
    void init();
    const onAuthChanged = () => setUser(getUser());
    const onPostStatsChanged = (e: Event) => {
      const { postId, likeCount, collectionCount: cc, commentCount: cmt } = (e as CustomEvent).detail;
      postsRef.current = postsRef.current.map(p =>
        p.id === postId
          ? { ...p, like_count: likeCount ?? p.like_count, favorite_count: cc ?? p.favorite_count, comment_count: cmt ?? p.comment_count }
          : p
      );
      setPosts(postsRef.current);
      // Update module cache too
      if (feedCache) {
        feedCache.posts = postsRef.current;
        feedCache.timestamp = Date.now();
      }
    };
    window.addEventListener("auth-changed", onAuthChanged);
    window.addEventListener("post-stats-changed", onPostStatsChanged);
    const onFeedRefresh = () => {
      feedCache = null;
      setTab("recommend");
      loadTab("recommend");
    };
    window.addEventListener("feed-refresh", onFeedRefresh);
    return () => {
      window.removeEventListener("auth-changed", onAuthChanged);
      window.removeEventListener("post-stats-changed", onPostStatsChanged);
      window.removeEventListener("feed-refresh", onFeedRefresh);
    };
  }, []);

  const switchTab = async (t: "recommend" | "following") => {
    if (t === tab) return;
    if (t === "following" && !user) {
      addToast("Please log in to view following feed.", { type: "warning", title: "Not logged in" });
      return;
    }
    setTab(t);
    feedCache = null; // clear cache on explicit tab switch
    await loadTab(t);
  };

  // --- Post creation ---
  const handleCreatePost = async () => {
    if (!user) {
      addToast("Please log in to create a post.", { type: "warning", title: "Not logged in" });
      return;
    }
    if (!title.trim() || !content.trim()) {
      addToast("Title and content are required.", { type: "warning", title: "Missing fields" });
      return;
    }
    setCreating(true);
    try {
      const response = await apiFetch("/api/posts", {
        method: "POST",
        body: JSON.stringify({ title: title.trim(), content: content.trim() }),
      });
      if (response.code !== 200) {
        addToast(response.message || "Failed to create post.", { type: "error", title: "Creation failed" });
        return;
      }
      const newPost = response.data;
      setPosts((current) => [{
        id: newPost.id, user_id: newPost.user_id, username: user.username,
        title: newPost.title, created_time: newPost.created_time,
        like_count: 0, comment_count: 0, favorite_count: 0, view_count: 0,
      }, ...current]);
      setTitle(""); setContent("");
      setShowCreateModal(false);
      feedCache = null; // invalidate cache so refresh fetches new post
      addToast("Post created successfully!", { type: "success", title: "Success" });
    } catch (error) {
      addToast("Unable to create post.", { type: "error", title: "Error" });
      console.error(error);
    } finally { setCreating(false); }
  };

  // --- infinite scroll observer ---
  const observerRef = useRef<IntersectionObserver | null>(null);
  const lastItemRef = useRef<HTMLDivElement | null>(null);
  const hasMoreRef = useRef(true);
  const loadingRef = useRef(false);
  const [observerPaused, setObserverPaused] = useState(false); // true during initial load / cache restore

  useEffect(() => { hasMoreRef.current = hasMore; }, [hasMore]);
  useEffect(() => { loadingRef.current = loading || loadingMore; }, [loading, loadingMore]);

  useEffect(() => {
    if (observerPaused) return;
    if (observerRef.current) observerRef.current.disconnect();
    observerRef.current = new IntersectionObserver((entries) => {
      if (observerPaused) return;
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        if (!loadingRef.current && hasMoreRef.current) {
          const endpoint = tab === "recommend" ? "/api/feed/recommend" : "/api/feed/following";
          void fetchFeed(endpoint, "subsequent", true);
        }
      }
    });
    const el = lastItemRef.current;
    if (el) observerRef.current.observe(el);
    return () => observerRef.current?.disconnect();
  }, [posts, tab, observerPaused]);

  return (
    <main className="mx-auto flex max-w-4xl flex-col gap-6 p-6">
      {/* ── Action bar: New Post + Refresh ── */}
      {user && (
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowCreateModal(true)}
            className="rounded bg-black px-4 py-2 text-sm text-white inline-flex items-center gap-1"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="8" y1="1" x2="8" y2="15" /><line x1="1" y1="8" x2="15" y2="8" />
            </svg>
            New Post
          </button>
          <button
            onClick={() => loadTab(tab)}
            className="rounded border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 inline-flex items-center gap-1"
          >
            <img src="/icon/refresh.svg" alt="refresh" style={{ width: 16, height: 16 }} />
            Refresh
          </button>
        </div>
      )}

      {/* ── Tab bar ── */}
      <div className="flex items-center gap-2 rounded border bg-gray-50 p-1">
        <button
          onClick={() => switchTab("recommend")}
          className={`flex-1 rounded px-4 py-2 text-sm font-medium transition ${
            tab === "recommend" ? "bg-white shadow text-black" : "text-gray-500 hover:text-gray-700"
          }`}
        >Recommend</button>
        <button
          onClick={() => switchTab("following")}
          className={`flex-1 rounded px-4 py-2 text-sm font-medium transition ${
            tab === "following" ? "bg-white shadow text-black" : "text-gray-500 hover:text-gray-700"
          }`}
        >Following</button>
      </div>

      {/* ── Post list ── */}
      <section className="space-y-4">
        {posts.length === 0 && loading ? (
          <div className="rounded border border-gray-200 bg-gray-50 p-4 text-gray-500">
            <Loading />
          </div>
        ) : posts.length > 0 ? (
          posts.map((post, idx) => {
            const observeIdx = posts.length - 1;
            return (
              <div key={post.id} ref={idx === observeIdx ? lastItemRef : undefined}>
                <PostCard post={post} />
              </div>
            );
          })
        ) : (
          <div className="rounded border border-gray-200 bg-gray-50 p-4 text-gray-500">
            No posts available.
          </div>
        )}

        {(loadingMore || loading) && posts.length > 0 && (
          <div className="py-6">
            <Loading />
          </div>
        )}

        {!hasMore && posts.length > 0 && !loadingMore && (
          <div className="rounded border border-gray-200 bg-gray-50 p-4 text-center text-sm text-gray-500">
            No more posts to load.
          </div>
        )}
      </section>

      {/* ── Create Post Modal ── */}
      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={() => setShowCreateModal(false)}>
          <div className="bg-white rounded-lg shadow-xl w-full max-w-lg p-6" onClick={(e) => e.stopPropagation()}>
            <h2 className="text-lg font-bold mb-4">Create Post</h2>
            <div className="space-y-3">
              <input
                className="w-full rounded border border-gray-300 bg-white p-3"
                placeholder="Title" value={title}
                onChange={(e) => setTitle(e.target.value)} disabled={creating}
              />
              <textarea
                className="w-full rounded border border-gray-300 bg-white p-3"
                placeholder="Content" rows={5} value={content}
                onChange={(e) => setContent(e.target.value)} disabled={creating}
              />
            </div>
            <div className="flex justify-end gap-2 mt-4">
              <button
                type="button"
                onClick={() => setShowCreateModal(false)}
                className="rounded border px-4 py-2 text-sm text-gray-600 hover:bg-gray-50"
              >Cancel</button>
              <button
                type="button" disabled={creating} onClick={handleCreatePost}
                className="rounded bg-black px-4 py-2 text-sm text-white hover:bg-gray-900 disabled:opacity-50"
              >{creating ? "Posting..." : "Publish"}</button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
