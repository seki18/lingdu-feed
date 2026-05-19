'use client';

import Link from "next/link";
import { useEffect, useState, use } from "react";
import { useSearchParams } from "next/navigation";
import { apiFetch, markPostDirty } from "@/lib/api";
import { useToast } from "@/components/ui/ToastContext";
import { PostDetail } from "@/types/post";
import { CommentItem } from "@/types/comment";
import { getUser } from "@/lib/auth";
import { User } from "@/types/user";
import { useRouter } from "next/navigation";
import CommentSection from "@/components/comment/CommentSection";

interface Props {
  params: Promise<{ id: string }>;
}

export default function PostDetailPage({ params }: Props) {
  const { id } = use(params);
  const router = useRouter();
  const searchParams = useSearchParams();
  const [post, setPost] = useState<PostDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const { addToast } = useToast();

  // Cached data from homepage (passed via searchParams)
  const cachedTitle = searchParams.get("t") || "";
  const cachedUsername = searchParams.get("u") || "";
  const cachedUserId = Number(searchParams.get("uid") ?? 0);
  const cachedCreatedTime = searchParams.get("ct") || "";
  const cachedPraiseCount = Number(searchParams.get("pc") ?? -1);
  const cachedCommentCount = Number(searchParams.get("cc") ?? -1);
  const cachedCollectionCount = Number(searchParams.get("clc") ?? -1);
  const cachedViewCount = Number(searchParams.get("vc") ?? -1);
  const cachedHasPraised = searchParams.get("hp") === "1";
  const cachedHasCollected = searchParams.get("hc") === "1";

  useEffect(() => {
    const u = getUser();
    if (u) setCurrentUser(u);
  }, []);

  useEffect(() => {
    const loadPost = async () => {
      setLoading(true);
      setNotFound(false);

      try {
        const response = await apiFetch(`/posts/${id}`);
        if (response.code !== 200) {
          if (response.code === 40004) {
            setNotFound(true);
          } else {
            addToast(response.message || "Unable to load post.", {
              type: "error",
              title: "Load failed",
            });
          }
          return;
        }
        const data = response.data;
        // New API: {post, has_praised, has_collected, comments}
        if (data?.post) {
          setPost(data.post);
          setHasPraised(data.has_praised ?? false);
          setHasCollected(data.has_collected ?? false);
          setPraiseCount(data.post.praise_count ?? 0);
          setCommentCount(data.post.comment_count ?? 0);
          setCollectionCount(data.post.collection_count ?? 0);
          setViewCount(data.post.view_count ?? 0);
          // Comments come from the API now, not CommentSection
          if (data.comments) setInitialComments(data.comments);
        } else {
          // Fallback for old API
          setPost(data);
        }
      } catch (error) {
        addToast("Unable to load post.", {
          type: "error",
          title: "Error",
        });
        console.error(error);
      } finally {
        setLoading(false);
      }
    };

    void loadPost();
  }, [id, addToast]);

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    setDeleting(true);
    try {
      const res = await apiFetch(`/posts/${id}`, { method: "DELETE" });
      if (res.code === 200) {
        addToast("Post deleted.", { type: "success", title: "Deleted" });
        router.push("/");
      } else {
        addToast(res.message || "Failed to delete.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setDeleting(false);
    }
  };

  // --- Stats (from homepage cache or single API call) ---
  const [praiseCount, setPraiseCount] = useState(cachedPraiseCount >= 0 ? cachedPraiseCount : 0);
  const [hasPraised, setHasPraised] = useState(cachedHasPraised);
  const [praising, setPraising] = useState(false);
  const [hasCollected, setHasCollected] = useState(cachedHasCollected);
  const [collecting, setCollecting] = useState(false);
  const [commentCount, setCommentCount] = useState(cachedCommentCount >= 0 ? cachedCommentCount : 0);
  const [collectionCount, setCollectionCount] = useState(cachedCollectionCount >= 0 ? cachedCollectionCount : 0);
  const [viewCount, setViewCount] = useState(cachedViewCount >= 0 ? cachedViewCount : 0);
  const [initialComments, setInitialComments] = useState<CommentItem[] | null>(null);

  // Stats come from the detail API — no separate /stats call needed

  const handlePraise = async () => {
    if (!currentUser) { addToast("Please log in to like.", { type: "warning", title: "Not logged in" }); return; }
    setPraising(true);
    try {
      if (hasPraised) {
        const res = await apiFetch("/Praises", { method: "DELETE", body: JSON.stringify({ post_id: Number(id) }) });
        if (res.code === 200) {
          setHasPraised(false);
          setPraiseCount((c) => Math.max(0, c - 1));
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), praiseCount: praiseCount - 1, hasPraised: false } }));
        }
      } else {
        const res = await apiFetch("/Praises", { method: "POST", body: JSON.stringify({ post_id: Number(id) }) });
        if (res.code === 200) {
          setHasPraised(true);
          setPraiseCount((c) => c + 1);
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), praiseCount: praiseCount + 1, hasPraised: true } }));
        }
      }
    } catch { addToast("Network error.", { type: "error", title: "Error" }); }
    finally { setPraising(false); }
  };

  const isAuthor = currentUser && post && currentUser.id === post.user_id;

  const handleCollect = async () => {
    if (!currentUser) { addToast("Please log in to collect.", { type: "warning", title: "Not logged in" }); return; }
    setCollecting(true);
    try {
      if (hasCollected) {
        const res = await apiFetch("/Collections", { method: "DELETE", body: JSON.stringify({ post_id: Number(id) }) });
        if (res.code === 200) {
          setHasCollected(false);
          setCollectionCount((c) => Math.max(0, c - 1));
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), collectionCount: collectionCount - 1, hasCollected: false } }));
        }
      } else {
        const res = await apiFetch("/Collections", { method: "POST", body: JSON.stringify({ post_id: Number(id) }) });
        if (res.code === 200) {
          setHasCollected(true);
          setCollectionCount((c) => c + 1);
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), collectionCount: collectionCount + 1, hasCollected: true } }));
        }
      }
    } catch { addToast("Network error.", { type: "error", title: "Error" }); }
    finally { setCollecting(false); }
  };

  return (
    <main className="mx-auto max-w-3xl p-6">
      <div className="mb-6 flex items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">Post Detail</h1>
          <p className="text-sm text-gray-600">Viewing post #{id}</p>
        </div>
        <div className="flex gap-2">
          <span className="rounded border px-4 py-2 text-sm inline-flex items-center gap-1 text-gray-500">
            <img src="/icon/see.svg" alt="view" style={{ width: 16, height: 16 }} />
            {viewCount}
          </span>
          <button
            type="button"
            onClick={handlePraise}
            disabled={praising}
            className={`rounded border px-4 py-2 text-sm inline-flex items-center gap-1 disabled:opacity-50 ${
              hasPraised ? "border-red-300 text-red-600 bg-red-50" : "border-gray-300 text-gray-700 hover:bg-gray-50"
            }`}
          >
            <img src={hasPraised ? "/icon/praise_yes.svg" : "/icon/praise_no.svg"} alt="praise" style={{ width: 16, height: 16 }} />
            {praiseCount}
          </button>
          <button
            type="button"
            onClick={handleCollect}
            disabled={collecting}
            className={`rounded border px-4 py-2 text-sm inline-flex items-center gap-1 disabled:opacity-50 ${
              hasCollected ? "border-yellow-300 text-yellow-600 bg-yellow-50" : "border-gray-300 text-gray-700 hover:bg-gray-50"
            }`}
          >
            <img src={hasCollected ? "/icon/collect_yes.svg" : "/icon/collect_no.svg"} alt="collect" style={{ width: 16, height: 16 }} />
            {hasCollected ? "Collected" : "Collect"}
          </button>
          {isAuthor && (
            <button
              type="button"
              onClick={handleDelete}
              disabled={deleting}
              className="rounded border border-red-300 bg-white px-4 py-2 text-sm text-red-600 hover:bg-red-50 disabled:opacity-50 inline-flex items-center gap-1"
            >
              <img src="/icon/delete.svg" alt="delete" style={{ width: 16, height: 16 }} />
              {deleting ? "Deleting..." : "Delete"}
            </button>
          )}
          <Link
            href="/"
            className="rounded border border-gray-300 bg-white px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 inline-flex items-center gap-1"
          >
            <img src="/icon/back.svg" alt="back" style={{ width: 16, height: 16 }} />
            Back to list
          </Link>
        </div>
      </div>

      {loading ? (
        <div className="rounded border border-gray-200 bg-gray-50 p-6 text-gray-500">
          Loading post...
        </div>
      ) : notFound ? (
        <div className="rounded border border-gray-200 bg-gray-50 p-6 text-gray-500">
          Post not found.
        </div>
      ) : post ? (
        <article className="space-y-6 rounded border border-gray-200 bg-white p-6 shadow-sm">
          <div className="space-y-2">
            <h2 className="text-3xl font-bold">{post.title}</h2>
            <div className="flex flex-wrap items-center gap-2 text-sm text-gray-500">
              <span>Post ID: {post.id}</span>
              <Link href={`/users/${post.user_id}`} className="text-blue-600 hover:underline">
                {post.username || `User #${post.user_id}`}
              </Link>
              <span>Created: {new Date(post.created_time).toLocaleString()}</span>
              <span>Updated: {new Date(post.updated_time).toLocaleString()}</span>
            </div>
          </div>
          <div className="whitespace-pre-line text-gray-700">{post.content}</div>
        </article>
      ) : (
        <div className="rounded border border-gray-200 bg-gray-50 p-6 text-gray-500">
          Post not found.
        </div>
      )}

      {post && (
        <div className="mt-8 pt-6 border-t border-gray-200">
          <CommentSection postId={id} initialComments={initialComments} initialCommentCount={commentCount} />
        </div>
      )}
    </main>
  );
}
