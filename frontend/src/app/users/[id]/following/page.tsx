'use client';

import { useEffect, useState, use } from "react";
import { getFollowingList, unfollowUser, type FollowItem } from "@/lib/api";
import { getUser } from "@/lib/auth";
import Link from "next/link";

interface Props {
  params: Promise<{ id: string }>;
}

export default function FollowingPage({ params }: Props) {
  const { id } = use(params);
  const [follows, setFollows] = useState<FollowItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [unfollowId, setUnfollowId] = useState<number | null>(null);
  const pageSize = 10;
  const currentUser = getUser();
  const isOwner = currentUser && Number(id) === currentUser.id;

  const fetchList = async (p = 1) => {
    setLoading(true);
    try {
      const res = await getFollowingList(Number(id), p, pageSize);
      if (res.code === 200) {
        setFollows(res.data?.follows ?? []);
        setTotal(res.data?.total ?? 0);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchList(page);
  }, [id, page]);

  const handleUnfollow = async (followingId: number) => {
    if (!confirm("Unfollow this user?")) return;
    setUnfollowId(followingId);
    try {
      const res = await unfollowUser(followingId);
      if (res.code === 200) {
        setFollows((prev) => prev.filter((f) => f.following_id !== followingId));
        setTotal((t) => t - 1);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setUnfollowId(null);
    }
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <main className="mx-auto max-w-2xl p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Following</h1>
        <Link href={`/users/${id}`} className="rounded border px-4 py-2 text-sm hover:bg-gray-50">
          Back to profile
        </Link>
      </div>

      {loading ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500 text-center">Loading...</div>
      ) : follows.length === 0 ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500 text-center">
          {isOwner ? "You are not following anyone yet." : "This user is not following anyone yet."}
        </div>
      ) : (
        <div className="space-y-2">
          {follows.map((item) => (
            <div key={item.following_id} className="flex items-center justify-between rounded border p-3">
              <Link
                href={`/users/${item.following_id}`}
                className="flex items-center gap-3 hover:text-blue-600"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src="/icon/me.svg" alt="user" style={{ width: 24, height: 24 }} />
                <span className="font-medium">{item.username || `User #${item.following_id}`}</span>
              </Link>
              {isOwner && (
                <button
                  onClick={() => handleUnfollow(item.following_id)}
                  disabled={unfollowId === item.following_id}
                  className="rounded border border-gray-300 px-3 py-1 text-sm text-gray-600 hover:bg-gray-100 disabled:opacity-50"
                >
                  {unfollowId === item.following_id ? "..." : "Unfollow"}
                </button>
              )}
            </div>
          ))}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 pt-4">
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
    </main>
  );
}
