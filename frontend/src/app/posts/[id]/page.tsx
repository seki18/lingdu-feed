'use client';

import Link from "next/link";
import { useEffect, useState, use } from "react";
import { apiFetch, markPostDirty, likePost, unlikePost, favoritePost, unfavoritePost } from "@/lib/api";
import { useToast } from "@/components/ui/ToastContext";
import { Post, PostDetailResponse } from "@/types/post";
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
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const { addToast } = useToast();

  // All data comes from the API — no more URL searchParams caching
  const [likeCount, setLikeCount] = useState(0);
  const [hasLiked, setHasLiked] = useState(false);
  const [liking, setLiking] = useState(false);
  const [hasFavorited, setHasFavorited] = useState(false);
  const [favoriting, setFavoriting] = useState(false);
  const [viewCount, setViewCount] = useState(0);
  const [initialComments, setInitialComments] = useState<CommentItem[] | null>(null);

  useEffect(() => {
    const u = getUser();
    if (u) setCurrentUser(u);
  }, []);

  useEffect(() => {
    const loadPost = async () => {
      setLoading(true);
      setNotFound(false);

      try {
        const response = await apiFetch<PostDetailResponse>(`/api/posts/${id}`);
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
        if (data?.post) {
          setPost(data.post);
          setHasLiked(data.has_liked ?? false);
          setHasFavorited(data.has_favorited ?? false);
          setLikeCount(data.post.like_count ?? 0);
          setViewCount(Math.max(data.post.view_count ?? 0, 1)); // optimistically +1 for this view
          setInitialComments(data.comments ?? []);
          // Only mark dirty on actual interactions (like/favorite/comment), not on page view
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
      const res = await apiFetch(`/api/posts/${id}`, { method: "DELETE" });
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

  const handleLike = async () => {
    if (!currentUser) { addToast("Please log in to like.", { type: "warning", title: "Not logged in" }); return; }
    setLiking(true);
    try {
      if (hasLiked) {
        const res = await unlikePost(Number(id));
        if (res.code === 200) {
          setHasLiked(false);
          setLikeCount((c) => Math.max(0, c - 1));
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), likeCount: likeCount - 1 } }));
        }
      } else {
        const res = await likePost(Number(id));
        if (res.code === 200) {
          setHasLiked(true);
          setLikeCount((c) => c + 1);
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), likeCount: likeCount + 1 } }));
        }
      }
    } catch { addToast("Network error.", { type: "error", title: "Error" }); }
    finally { setLiking(false); }
  };

  const handleFavorite = async () => {
    if (!currentUser) { addToast("Please log in to favorite.", { type: "warning", title: "Not logged in" }); return; }
    setFavoriting(true);
    try {
      if (hasFavorited) {
        const res = await unfavoritePost(Number(id));
        if (res.code === 200) {
          setHasFavorited(false);
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), hasFavorited: false } }));
        }
      } else {
        const res = await favoritePost(Number(id));
        if (res.code === 200) {
          setHasFavorited(true);
          markPostDirty(Number(id));
          window.dispatchEvent(new CustomEvent("post-stats-changed", { detail: { postId: Number(id), hasFavorited: true } }));
        }
      }
    } catch { addToast("Network error.", { type: "error", title: "Error" }); }
    finally { setFavoriting(false); }
  };

  const isAuthor = currentUser && post && currentUser.id === post.user_id;

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
            onClick={handleLike}
            disabled={liking}
            className={`rounded border px-4 py-2 text-sm inline-flex items-center gap-1 disabled:opacity-50 ${
              hasLiked ? "border-red-300 text-red-600 bg-red-50" : "border-gray-300 text-gray-700 hover:bg-gray-50"
            }`}
          >
            <img src={hasLiked ? "/icon/praise_yes.svg" : "/icon/praise_no.svg"} alt="like" style={{ width: 16, height: 16 }} />
            {likeCount}
          </button>
          <button
            type="button"
            onClick={handleFavorite}
            disabled={favoriting}
            className={`rounded border px-4 py-2 text-sm inline-flex items-center gap-1 disabled:opacity-50 ${
              hasFavorited ? "border-yellow-300 text-yellow-600 bg-yellow-50" : "border-gray-300 text-gray-700 hover:bg-gray-50"
            }`}
          >
            <img src={hasFavorited ? "/icon/collect_yes.svg" : "/icon/collect_no.svg"} alt="favorite" style={{ width: 16, height: 16 }} />
            {hasFavorited ? "Favorited" : "Favorite"}
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
          <CommentSection postId={id} initialComments={initialComments} />
        </div>
      )}
    </main>
  );
}
