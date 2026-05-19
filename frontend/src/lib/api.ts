import { getToken } from "./auth";

const BASE_URL = "http://localhost:18080";

// ── Dirty post tracker (module-level, survives component remounts) ──
// Detail page sets dirtyPostId when user interacts (like/collect/comment).
// Feed page reads & clears it after restoring from cache.
let dirtyPostIds = new Set<number>();
export function markPostDirty(postId: number): void {
  dirtyPostIds.add(postId);
}
export function consumeDirtyPostIds(): number[] {
  const ids = Array.from(dirtyPostIds);
  dirtyPostIds.clear();
  return ids;
}

// ApiResponse is the standard JSON envelope returned by all backend endpoints.
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// apiFetch sends an authenticated request to the backend API and returns the parsed response.
// Pass skipAuth: true in options to omit the Authorization header.
export async function apiFetch<T = any>(
  path: string,
  options: RequestInit & { skipAuth?: boolean } = {}
): Promise<ApiResponse<T>> {
  const { skipAuth, ...fetchOptions } = options;
  const token = getToken();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token && !skipAuth) {
    headers.Authorization = `Bearer ${token}`;
  }

  try {
    const response = await fetch(`${BASE_URL}${path}`, {
      ...fetchOptions,
      headers,
    });

    const data: ApiResponse<T> = await response.json();
    return data;
  } catch (error) {
    return {
      code: 50000,
      message: error instanceof Error ? error.message : "An unexpected network error occurred",
    };
  }
}

// ── Interaction status (with client-side cache to avoid redundant requests) ──

const statusCache = new Map<number, number>(); // postId → highest status already sent
let pendingBatch: { post_id: number; status: 1 | 2 | 3 }[] = [];
let batchTimer: ReturnType<typeof setTimeout> | null = null;

function flushBatch() {
  if (pendingBatch.length === 0) return;
  const batch = pendingBatch;
  pendingBatch = [];
  apiFetch("/interaction-status/batch", {
    method: "POST",
    body: JSON.stringify(batch),
  }).catch((err) => console.warn("Batch status update failed:", err));
}

// trackInteractionStatus tracks user interactions with posts (delivery=1, display=2, click=3).
// Skips if the new status is not higher than what was already recorded locally,
// EXCEPT for status=3 (click) which is always sent to update updated_time for history ordering.
export function trackInteractionStatus(postId: number, status: 1 | 2 | 3): void {
  const prev = statusCache.get(postId) ?? 0;
  if (status < prev) return;
  if (status === prev && status !== 3) return;
  statusCache.set(postId, status);

  pendingBatch.push({ post_id: postId, status });
  if (!batchTimer) {
    batchTimer = setTimeout(() => {
      batchTimer = null;
      flushBatch();
    }, 500);
  }
}

// batchTrackDelivery immediately reports a list of post IDs as delivered (status=1).
export function batchTrackDelivery(postIds: number[]): void {
  for (const id of postIds) {
    const prev = statusCache.get(id) ?? 0;
    if (1 <= prev) continue;
    statusCache.set(id, 1);
  }
  const batch = postIds
    .filter((id) => !statusCache.has(id) || statusCache.get(id)! === 1)
    .map((id) => ({ post_id: id, status: 1 as const }));
  if (batch.length === 0) return;
  apiFetch("/interaction-status/batch", {
    method: "POST",
    body: JSON.stringify(batch),
  }).catch((err) => console.warn("Batch delivery update failed:", err));
}

// ── Post stats (bundled exist + counts for detail page) ──

// updateUsername sends a PUT /users request to change the current user's username (auth required).
export async function updateUsername(username: string): Promise<ApiResponse> {
  return apiFetch("/users", {
    method: "PUT",
    body: JSON.stringify({ username }),
  });
}

// changePassword sends a PUT /users/password request to change the current user's password (auth required).
export async function changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse> {
  return apiFetch("/users/password", {
    method: "PUT",
    body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
  });
}

// followUser follows a user (auth required).
export async function followUser(followingId: number): Promise<ApiResponse> {
  return apiFetch("/Follows", {
    method: "POST",
    body: JSON.stringify({ following_id: followingId }),
  });
}

// unfollowUser unfollows a user (auth required).
export async function unfollowUser(followingId: number): Promise<ApiResponse> {
  return apiFetch("/Follows", {
    method: "DELETE",
    body: JSON.stringify({ following_id: followingId }),
  });
}

export interface FollowItem {
  follower_id: number;
  following_id: number;
  username: string;
  created_time: string;
}

// getFollowingList returns the list of users that userId is following (paginated).
export async function getFollowingList(
  userId: number,
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ follows: FollowItem[]; total: number }>> {
  return apiFetch(`/Follows/list/following/${userId}?page=${page}&pageSize=${pageSize}`);
}

// getFollowerList returns the list of followers of userId (paginated).
export async function getFollowerList(
  userId: number,
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ follows: FollowItem[]; total: number }>> {
  return apiFetch(`/Follows/list/follower/${userId}?page=${page}&pageSize=${pageSize}`);
}