'use client';

import { useEffect, useState, useRef } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import { useToast } from "@/components/ui/ToastContext";
import { CommentItem } from "@/types/comment";
import { getUser } from "@/lib/auth";

interface Props {
  postId: string;
  initialComments?: CommentItem[] | null;
  initialCommentCount?: number;
}

// Inline SVG icon helper – uses SVG files from /public/icon/
function Icon({ name, className }: { name: string; className?: string }) {
  return (
    <img src={`/icon/${name}.svg`} alt={name} className={className} style={{ width: 16, height: 16, display: "inline" }} />
  );
}

export default function CommentSection({ postId, initialComments, initialCommentCount }: Props) {
  const [comments, setComments] = useState<CommentItem[]>(initialComments ?? []);
  const [loading, setLoading] = useState(!initialComments);
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [replyTo, setReplyTo] = useState<number | null>(null);
  const [replyContent, setReplyContent] = useState("");
  const [commentCount, setCommentCount] = useState(initialCommentCount ?? 0);
  const [currentUserId, setCurrentUserId] = useState<number | null>(null);
  const initialLoaded = useRef(!!initialComments);
  const { addToast } = useToast();

  useEffect(() => {
    const u = getUser();
    setCurrentUserId(u?.id ?? null);
  }, []);

  const fetchComments = async () => {
    if (initialLoaded.current) return;
    setLoading(true);
    try {
      const [commentRes, countRes] = await Promise.all([
        apiFetch(`/comments/by-post/${postId}`),
        apiFetch(`/comments/count/${postId}`),
      ]);
      if (commentRes.code === 200) {
        setComments(commentRes.data ?? []);
      }
      if (countRes.code === 200) {
        setCommentCount(countRes.data?.count ?? 0);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchComments();
  }, [postId]);

  const handleSubmit = async () => {
    if (!content.trim()) {
      addToast("Comment cannot be empty.", { type: "warning", title: "Missing content" });
      return;
    }
    setSubmitting(true);
    try {
      const res = await apiFetch("/comments", {
        method: "POST",
        body: JSON.stringify({ post_id: Number(postId), content: content.trim() }),
      });
      if (res.code === 200) {
        setContent("");
        addToast("Comment added.", { type: "success", title: "Success" });
        await fetchComments();
      } else {
        addToast(res.message || "Failed to add comment.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setSubmitting(false);
    }
  };

  const handleReply = async () => {
    if (!replyContent.trim() || replyTo === null) return;
    setSubmitting(true);
    try {
      const res = await apiFetch("/comments", {
        method: "POST",
        body: JSON.stringify({
          post_id: Number(postId),
          content: replyContent.trim(),
          reply_id: replyTo,
        }),
      });
      if (res.code === 200) {
        setReplyTo(null);
        setReplyContent("");
        addToast("Reply added.", { type: "success", title: "Success" });
        await fetchComments();
      } else {
        addToast(res.message || "Failed to reply.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (commentId: number) => {
    if (!confirm("Delete this comment? Its replies will also be removed.")) return;
    try {
      const res = await apiFetch(`/comments/${commentId}`, { method: "DELETE" });
      if (res.code === 200) {
        addToast("Comment deleted.", { type: "success", title: "Deleted" });
        await fetchComments();
      } else {
        addToast(res.message || "Failed to delete.", { type: "error", title: "Error" });
      }
    } catch {
      addToast("Network error.", { type: "error", title: "Error" });
    }
  };

  return (
    <section className="space-y-4">
      <h2 className="text-xl font-bold inline-flex items-center gap-1">
        <Icon name="comment" />
        Comments ({commentCount})
      </h2>

      {/* New comment input */}
      <div className="space-y-2">
        <textarea
          className="w-full rounded border border-gray-300 p-3 text-sm"
          rows={3}
          placeholder="Write a comment..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          disabled={submitting}
        />
        <button
          type="button"
          onClick={handleSubmit}
          disabled={submitting}
          className="rounded bg-black px-4 py-2 text-sm text-white hover:bg-gray-800 disabled:opacity-50"
        >
          {submitting ? "Submitting..." : "Comment"}
        </button>
      </div>

      {/* Two-layer grouping */}
      {loading ? (
        <div className="rounded border border-gray-200 bg-gray-50 p-4 text-sm text-gray-500">
          Loading comments...
        </div>
      ) : comments.length === 0 ? (
        <div className="rounded border border-gray-200 bg-gray-50 p-4 text-sm text-gray-500">
          No comments yet. Be the first to comment!
        </div>
      ) : (
        <div className="space-y-3">
          {(() => {
            const topLevel = comments.filter((c) => !c.reply_id);

            const getReplies = (rootId: number): CommentItem[] =>
              comments.filter((c) => {
                if (!c.reply_id) return false;
                let parentId: number | null = c.reply_id;
                while (parentId !== null) {
                  if (parentId === rootId) return true;
                  const parent = comments.find((p) => p.id === parentId);
                  parentId = parent?.reply_id ?? null;
                }
                return false;
              });

            return topLevel.map((root) => {
              const replies = getReplies(root.id);
              return (
                <div key={root.id} className="rounded border border-gray-200 p-3">
                  {/* Top-level comment */}
                  <div className="flex items-center gap-2 mb-1">
                    <Link href={`/users/${root.user_id}`} className="text-sm font-medium hover:text-blue-600 hover:underline">{root.username}</Link>
                    <span className="text-xs text-gray-400">
                      {new Date(root.created_time).toLocaleString()}
                    </span>
                  </div>
                  <p className="text-sm text-gray-700">{root.content}</p>

                  {/* Action bar */}
                  <div className="mt-2 flex items-center gap-3">
                    <button
                      type="button"
                      onClick={() => setReplyTo(root.id)}
                      className="text-xs text-gray-500 hover:text-gray-700 inline-flex items-center gap-1"
                    >
                      <Icon name="comment" /> Reply
                    </button>
                    {currentUserId === root.user_id && (
                      <button
                        type="button"
                        onClick={() => handleDelete(root.id)}
                        className="text-xs text-red-500 hover:text-red-700 inline-flex items-center gap-1"
                      >
                        <Icon name="delete" /> Delete
                      </button>
                    )}
                  </div>

                  {/* Reply form */}
                  {replyTo === root.id && (
                    <div className="mt-2 space-y-2 pl-4 border-l-2 border-gray-200">
                      <textarea
                        className="w-full rounded border border-gray-300 p-2 text-sm"
                        rows={2}
                        placeholder={`Reply to @${root.username}...`}
                        value={replyContent}
                        onChange={(e) => setReplyContent(e.target.value)}
                        disabled={submitting}
                      />
                      <div className="flex gap-2">
                        <button
                          type="button"
                          onClick={handleReply}
                          disabled={submitting}
                          className="rounded bg-black px-3 py-1 text-xs text-white hover:bg-gray-800 disabled:opacity-50"
                        >
                          {submitting ? "..." : "Reply"}
                        </button>
                        <button
                          type="button"
                          onClick={() => { setReplyTo(null); setReplyContent(""); }}
                          className="rounded border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  )}

                  {/* Nested replies */}
                  {replies.length > 0 && (
                    <div className="mt-2 space-y-2">
                      {replies.map((reply) => (
                        <div key={reply.id} className="pl-6 border-l-2 border-gray-200">
                          <div className="flex items-center gap-2 mb-1">
                            <Link href={`/users/${reply.user_id}`} className="text-sm font-medium hover:text-blue-600 hover:underline">{reply.username}</Link>
                            {reply.reply_username && (
                              <span className="text-xs text-blue-500">
                                @{reply.reply_username}
                              </span>
                            )}
                            <span className="text-xs text-gray-400">
                              {new Date(reply.created_time).toLocaleString()}
                            </span>
                          </div>
                          <p className="text-sm text-gray-700">{reply.content}</p>

                          {/* Action bar for replies */}
                          <div className="mt-2 flex items-center gap-3">
                            <button
                              type="button"
                              onClick={() => setReplyTo(reply.id)}
                              className="text-xs text-gray-500 hover:text-gray-700 inline-flex items-center gap-1"
                            >
                              <Icon name="comment" /> Reply
                            </button>
                            {currentUserId === reply.user_id && (
                              <button
                                type="button"
                                onClick={() => handleDelete(reply.id)}
                                className="text-xs text-red-500 hover:text-red-700 inline-flex items-center gap-1"
                              >
                                <Icon name="delete" /> Delete
                              </button>
                            )}
                          </div>

                          {/* Reply form for nested */}
                          {replyTo === reply.id && (
                            <div className="mt-2 space-y-2">
                              <textarea
                                className="w-full rounded border border-gray-300 p-2 text-sm"
                                rows={2}
                                placeholder={`Reply to @${reply.username}...`}
                                value={replyContent}
                                onChange={(e) => setReplyContent(e.target.value)}
                                disabled={submitting}
                              />
                              <div className="flex gap-2">
                                <button
                                  type="button"
                                  onClick={handleReply}
                                  disabled={submitting}
                                  className="rounded bg-black px-3 py-1 text-xs text-white hover:bg-gray-800 disabled:opacity-50"
                                >
                                  {submitting ? "..." : "Reply"}
                                </button>
                                <button
                                  type="button"
                                  onClick={() => { setReplyTo(null); setReplyContent(""); }}
                                  className="rounded border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
                                >
                                  Cancel
                                </button>
                              </div>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              );
            });
          })()}
        </div>
      )}
    </section>
  );
}
