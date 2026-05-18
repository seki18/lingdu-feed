'use client';

import { useEffect, useState } from "react";
import PostCard from "@/components/layout/PostCard";
import { apiFetch } from "@/lib/api";
import { useToast } from "@/components/ui/ToastContext";
import { PostSummary } from "@/types/post";

export default function CollectionsPage() {
  const [items, setItems] = useState<PostSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 10;
  const { addToast } = useToast();

  useEffect(() => {
    const fetchData = async (p: number) => {
      setLoading(true);
      try {
        const res = await apiFetch(`/feed/collections?page=${p}&page_size=${pageSize}`);
        if (res.code === 200) {
          setItems(res.data?.items ?? []);
          setTotal(res.data?.total ?? 0);
        } else {
          addToast(res.message || "Failed to load collections.", { type: "error", title: "Error" });
        }
      } catch {
        addToast("Network error.", { type: "error", title: "Error" });
      } finally {
        setLoading(false);
      }
    };
    void fetchData(page);
  }, [page, addToast]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <main className="mx-auto max-w-4xl p-6">
      <h1 className="text-2xl font-bold mb-4 inline-flex items-center gap-2">
        <img src="/icon/collect_yes.svg" alt="collection" style={{ width: 24, height: 24 }} />
        My Collections
      </h1>

      {loading ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500">Loading...</div>
      ) : items.length === 0 ? (
        <div className="rounded border bg-gray-50 p-6 text-gray-500">No collections yet.</div>
      ) : (
        <div className="space-y-3">
          {items.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center gap-2 mt-4">
          <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1} className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-100">Prev</button>
          <span className="px-3 py-1 text-sm text-gray-600">{page} / {totalPages}</span>
          <button onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page === totalPages} className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-100">Next</button>
        </div>
      )}
    </main>
  );
}
