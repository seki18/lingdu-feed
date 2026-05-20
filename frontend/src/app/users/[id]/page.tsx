'use client';

import { useEffect, useState, use } from "react";
import { getUserFeed, getPostDetail, updatePost, deletePost, trackState, followUser, unfollowUser } from "@/lib/api";
import { getUser } from "@/lib/auth";
import { useToast } from "@/components/ui/ToastContext";
import { PostSummary } from "@/types/post";
import { User, UserProfilePage } from "@/types/user";
import Link from "next/link";
import ProfileModal from "@/components/auth/ProfileModal";

interface Props {
  params: Promise<{ id: string }>;
}

export default function UserPage({ params }: Props) {
  const { id } = use(params);
  const [user, setUser] = useState<User | null>(null);
  const [posts, setPosts] = useState<PostSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 10;
  const [editPostId, setEditPostId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState<number | null>(null);
  const [following, setFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [followingCount, setFollowingCount] = useState(0);
  const [followerCount, setFollowerCount] = useState(0);
  const [profileModalOpen, setProfileModalOpen] = useState(false);
  const currentUser = getUser();
  const isOwner = currentUser && Number(id) === currentUser.id;
  const { addToast } = useToast();

  const fetchUserAndPosts = async (p = 1) => {
    setLoading(true);
    try {
      const userRes = await getUserFeed(Number(id), p, pageSize);
      if (userRes.code === 200 && userRes.data) {
        setUser(userRes.data.user);
        setPosts(userRes.data.posts ?? []);
        setTotal(userRes.data.total ?? 0);
        setFollowing(userRes.data.user?.is_following ?? false);
        setFollowingCount(userRes.data.user?.following_count ?? 0);
        setFollowerCount(userRes.data.user?.follower_count ?? 0);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchUserAndPosts(page);
  }, [id, page]);

  const totalPages = Math.ceil(total / pageSize);

  const handleEdit = async (post: PostSummary) => {
    setEditPostId(post.id);
    setEditTitle(post.title);
    setEditContent("");
    // Fetch full post detail to get the content for editing
    try {
      const res = await getPostDetail(post.id);
      if (res.code === 200 && res.data?.content) {
        setEditContent(res.data.content);
      }
    } catch {
      // Leave content empty if fetch fails
    }
  };

  const handleCancelEdit = () => {
    setEditPostId(null);
    setEditTitle("");
    setEditContent("");
  };

  const handleSaveEdit = async () => {
    if (!editTitle.trim() || !editContent.trim()) {
      addToast("Title and content are required.", { type: "warning", title: "Missing fields" });
      return;
    }
    setSaving(true);
    try {
      if (editPostId == null) return;
      const res = await updatePost(editPostId, editTitle.trim(), editContent.trim());
      if (res.code === 200) {
        addToast("Post updated successfully!", { type: "success", title: "Success" });
        handleCancelEdit();
        await fetchUserAndPosts();
      } else {
        addToast(res.message || "Failed to update.", { type: "error", title: "Update failed" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setSaving(false);
    }
  };

  const handleDeletePost = async (postId: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    setDeleting(postId);
    try {
      const res = await deletePost(postId);
      if (res.code === 200) {
        addToast("Post deleted.", { type: "success", title: "Deleted" });
        await fetchUserAndPosts();
      } else {
        addToast(res.message || "Failed to delete.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setDeleting(null);
    }
  };

  const handleFollow = async () => {
    if (!currentUser) {
      addToast("Please log in to follow.", { type: "warning", title: "Not logged in" });
      return;
    }
    setFollowLoading(true);
    try {
      if (following) {
        const res = await unfollowUser(Number(id));
        if (res.code === 200) {
          setFollowing(false);
          setFollowerCount((c) => c - 1);
        }
      } else {
        const res = await followUser(Number(id));
        if (res.code === 200) {
          setFollowing(true);
          setFollowerCount((c) => c + 1);
        }
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setFollowLoading(false);
    }
  };

  if (loading) {
    return (
      <main className="mx-auto max-w-4xl p-6">
        <div className="rounded border border-gray-200 bg-gray-50 p-6 text-gray-500">Loading...</div>
      </main>
    );
  }

  return (
    <main className="mx-auto max-w-4xl p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{user?.username || `User #${id}`}</h1>
          <p className="text-sm text-gray-600">{user?.email}</p>
          <div className="flex gap-4 mt-2 text-sm text-gray-500">
            <Link
              href={`/users/${id}/following`}
              className="hover:text-blue-600 hover:underline"
            >{followingCount} Following</Link>
            <Link
              href={`/users/${id}/followers`}
              className="hover:text-blue-600 hover:underline"
            >{followerCount} Followers</Link>
          </div>
        </div>
        <div className="flex items-center gap-2">
          {isOwner && (
            <button
              className="rounded border px-4 py-2 text-sm hover:bg-gray-100 inline-flex items-center gap-1"
              onClick={() => setProfileModalOpen(true)}
            >
              <img src="/icon/edit.svg" alt="edit" style={{ width: 16, height: 16 }} />
              Edit Profile
            </button>
          )}
          {!isOwner && currentUser && (
            <button
              onClick={handleFollow}
              disabled={followLoading}
              className={`rounded border px-4 py-2 text-sm font-medium disabled:opacity-50 ${
                following
                  ? "border-gray-300 text-gray-700 bg-gray-100 hover:bg-gray-200"
                  : "border-blue-500 bg-blue-500 text-white hover:bg-blue-600"
              }`}
            >
              {followLoading ? "..." : following ? "Unfollow" : "Follow"}
            </button>
          )}
          <Link href="/" className="rounded border px-4 py-2 text-sm hover:bg-gray-50 inline-flex items-center gap-1">
            <img src="/icon/back.svg" alt="back" style={{ width: 16, height: 16 }} />
            Back
          </Link>
        </div>
      </div>

      {/* Profile Modal */}
      {user && (
        <ProfileModal
          open={profileModalOpen}
          onClose={() => setProfileModalOpen(false)}
          user={user}
          onUserUpdated={(updated) => setUser(updated)}
        />
      )}

      <section>
        <h2 className="mb-4 text-lg font-semibold">Posts ({posts.length})</h2>
        {posts.length === 0 ? (
          <div className="rounded border bg-gray-50 p-6 text-gray-500">No posts yet.</div>
        ) : (
          <div className="space-y-3">
            {posts.map((post) => (
              <div key={post.id} className="rounded border p-4 shadow-sm">
                {editPostId === post.id ? (
                  <div className="space-y-3">
                    <input
                      className="w-full rounded border p-2"
                      value={editTitle}
                      onChange={(e) => setEditTitle(e.target.value)}
                      disabled={saving}
                      placeholder="Title"
                    />
                    <textarea
                      className="w-full rounded border p-2"
                      rows={3}
                      value={editContent}
                      onChange={(e) => setEditContent(e.target.value)}
                      disabled={saving}
                      placeholder="Content"
                    />
                    <div className="flex gap-2">
                      <button
                        className="rounded bg-black px-4 py-2 text-white disabled:opacity-50"
                        onClick={handleSaveEdit}
                        disabled={saving}
                      >
                        {saving ? "Saving..." : "Save"}
                      </button>
                      <button
                        className="rounded border px-4 py-2 disabled:opacity-50"
                        onClick={handleCancelEdit}
                        disabled={saving}
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <>
                    <div className="flex items-center justify-between gap-2">
                      <Link
                        href={`/posts/${post.id}`}
                        className="text-lg font-bold hover:underline"
                        onClick={() => { trackState(post.id, 3); }}
                      >
                        {post.title}
                      </Link>
                      <div className="flex items-center gap-3 text-sm text-gray-500">
                        <span className="inline-flex items-center gap-1">
                          <img src="/icon/see.svg" alt="" style={{ width: 14, height: 14 }} /> {post.view_count ?? 0}
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <img src="/icon/praise_no.svg" alt="" style={{ width: 14, height: 14 }} /> {post.like_count ?? 0}
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <img src="/icon/comment.svg" alt="" style={{ width: 14, height: 14 }} /> {post.comment_count ?? 0}
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <img src="/icon/collect_no.svg" alt="" style={{ width: 14, height: 14 }} /> {post.favorite_count ?? 0}
                        </span>
                        {isOwner && (
                          <div className="flex gap-1 ml-2">
                            <button onClick={() => handleEdit(post)} className="rounded border px-3 py-1 text-sm hover:bg-gray-100 inline-flex items-center gap-1">
                              <img src="/icon/edit.svg" alt="edit" style={{ width: 14, height: 14 }} />Edit
                            </button>
                            <button onClick={() => handleDeletePost(post.id)} disabled={deleting === post.id} className="rounded border border-red-300 px-3 py-1 text-sm text-red-600 hover:bg-red-50 disabled:opacity-50 inline-flex items-center gap-1">
                              <img src="/icon/delete.svg" alt="delete" style={{ width: 14, height: 14 }} />{deleting === post.id ? "..." : "Delete"}
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                    <p className="mt-1 text-sm text-gray-500">{new Date(post.created_time).toLocaleString()}</p>
                  </>
                )}
              </div>
            ))}
            {totalPages > 1 && (
              <div className="flex justify-center gap-2 mt-4">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-100"
                >
                  Prev
                </button>
                <span className="px-3 py-1 text-sm text-gray-600">{page} / {totalPages}</span>
                <button
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                  className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-100"
                >
                  Next
                </button>
              </div>
            )}
          </div>
        )}
      </section>
    </main>
  );
}