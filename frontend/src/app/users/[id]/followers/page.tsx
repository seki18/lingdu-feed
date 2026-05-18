'use client';

import { useEffect, useState, use } from "react";
import { getFollowerList, type FollowItem } from "@/lib/api";
import Link from "next/link";

interface Props {
  params: Promise<{ id: string }>;
}

export default function FollowersPage({ params }: Props) {
  const { id } = use(params);
  const [follows, setFollows] = useState<FollowItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 10;

  const fetchList = async (p = 1) => {
    setLoading(true);
    try {
      const res = await getFollowerList(Number(id), p, pageSize);
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

  const totalPages = Math.ceil(total / pageSize);

  return (
    <main className="mx-auto max-w-2xl p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Followers</h1>
        <Link href={`/users/${id}`} className="rounded border px-4 py-2 text-sm hover:bg-gray-50">
          Back to profile
        </Link>
      </div>

      {loading ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500 text-center">Loading...</div>
      ) : follows.length === 0 ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500 text-center">
          No followers yet.
        </div>
      ) : (
        <div className="space-y-2">
          {follows.map((item) => (
            <div key={item.follower_id} className="flex items-center justify-between rounded border p-3">
              <Link
                href={`/users/${item.follower_id}`}
                className="flex items-center gap-3 hover:text-blue-600"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src="/icon/me.svg" alt="user" style={{ width: 24, height: 24 }} />
                <span className="font-medium">{item.username || `User #${item.follower_id}`}</span>
              </Link>
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
