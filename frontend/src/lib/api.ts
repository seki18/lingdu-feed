import { getToken } from "./auth";
import { PostSummary } from "@/types/post";

const BASE_URL = "http://localhost:18080";

// ── Dirty post tracker (module-level, survives component remounts) ──
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

// ── State tracking (delivered=1, exposed=2, clicked=3) ──

const statusCache = new Map<number, number>();
let pendingBatch: { post_id: number; status: 1 | 2 | 3 }[] = [];
let batchTimer: ReturnType<typeof setTimeout> | null = null;

function flushBatch() {
  if (pendingBatch.length === 0) return;
  const batch = pendingBatch;
  pendingBatch = [];
  apiFetch("/api/state/batch", {
    method: "POST",
    body: JSON.stringify(batch),
  }).catch((err) => console.warn("Batch state update failed:", err));
}

export function trackState(postId: number, status: 1 | 2 | 3): void {
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

// ── Social API helpers ──

export async function likePost(postId: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/like`, { method: "POST" });
}

export async function unlikePost(postId: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/like`, { method: "DELETE" });
}

export async function favoritePost(postId: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/favorite`, { method: "POST" });
}

export async function unfavoritePost(postId: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/favorite`, { method: "DELETE" });
}

export async function getComments(postId: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/comments`);
}

export async function createComment(postId: number, content: string, replyId?: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/comments`, {
    method: "POST",
    body: JSON.stringify({ content, reply_id: replyId ?? null }),
  });
}

export async function deleteComment(commentId: number): Promise<ApiResponse> {
  return apiFetch(`/api/comments/${commentId}`, { method: "DELETE" });
}

// ── User API helpers ──

export async function updateUsername(username: string): Promise<ApiResponse> {
  return apiFetch("/api/users/me/profile", {
    method: "PUT",
    body: JSON.stringify({ username }),
  });
}

export async function changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse> {
  return apiFetch("/api/users/me/password", {
    method: "PUT",
    body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
  });
}

export async function followUser(userId: number): Promise<ApiResponse> {
  return apiFetch(`/api/users/${userId}/follow`, { method: "POST" });
}

export async function unfollowUser(userId: number): Promise<ApiResponse> {
  return apiFetch(`/api/users/${userId}/follow`, { method: "DELETE" });
}

export interface FollowItem {
  follower_id: number;
  following_id: number;
  username: string;
  created_time: string;
}

export async function getFollowingList(
  userId: number,
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ follows: FollowItem[]; total: number }>> {
  return apiFetch(`/api/users/${userId}/following?page=${page}&pageSize=${pageSize}`);
}

export async function getFollowerList(
  userId: number,
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ follows: FollowItem[]; total: number }>> {
  return apiFetch(`/api/users/${userId}/followers?page=${page}&pageSize=${pageSize}`);
}

// ── Feed page helpers ──

export async function getRecommendFeed(
  requestType: "initial" | "subsequent" = "subsequent",
  currentIds?: string
): Promise<ApiResponse<PostSummary[]>> {
  const qs = [`request_type=${requestType}`];
  if (currentIds) qs.push(`current_ids=${encodeURIComponent(currentIds)}`);
  return apiFetch(`/api/feed/recommend?${qs.join("&")}`);
}

export async function getFollowingFeed(
  requestType: "initial" | "subsequent" = "subsequent",
  currentIds?: string
): Promise<ApiResponse<PostSummary[]>> {
  const qs = [`request_type=${requestType}`];
  if (currentIds) qs.push(`current_ids=${encodeURIComponent(currentIds)}`);
  return apiFetch(`/api/feed/following?${qs.join("&")}`);
}

export async function getHistoryFeed(
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ items: PostSummary[]; total: number }>> {
  return apiFetch(`/api/feed/history?page=${page}&page_size=${pageSize}`);
}

export async function getFavoriteFeed(
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ items: PostSummary[]; total: number }>> {
  return apiFetch(`/api/feed/favorites?page=${page}&page_size=${pageSize}`);
}

export async function getUserFeed(
  userId: number,
  page = 1,
  pageSize = 10
): Promise<ApiResponse<{ user: any; posts: PostSummary[]; total: number }>> {
  return apiFetch(`/api/feed/users/${userId}?page=${page}&page_size=${pageSize}`);
}

export async function getPostDetail(id: number): Promise<ApiResponse<any>> {
  return apiFetch(`/api/posts/${id}`);
}

export async function updatePost(id: number, title: string, content: string): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${id}`, {
    method: "PUT",
    body: JSON.stringify({ title, content }),
  });
}

export async function deletePost(id: number): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${id}`, { method: "DELETE" });
}

export async function createPost(title: string, content: string): Promise<ApiResponse> {
  return apiFetch("/api/posts", {
    method: "POST",
    body: JSON.stringify({ title, content }),
  });
}

// addPostImages associates uploaded image URLs with a post.
export async function addPostImages(postId: number, images: string[]): Promise<ApiResponse> {
  return apiFetch(`/api/posts/${postId}/images`, {
    method: "POST",
    body: JSON.stringify({ images }),
  });
}

// uploadImage uploads a single file to S3 via the backend.
export async function uploadImage(postId: number, file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("post_id", String(postId));
  const token = getToken();
  const res = await fetch(`${BASE_URL}/api/upload`, {
    method: "POST",
    body: formData,
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  const data: ApiResponse = await res.json();
  if (data.code !== 200) throw new Error(data.message || "Upload failed");
  return data.data.url ?? data.data;
}

// ── Auth API helpers ──

export async function login(email: string, password: string): Promise<ApiResponse> {
  return apiFetch("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function register(email: string, password: string, username: string): Promise<ApiResponse> {
  return apiFetch("/api/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, username }),
  });
}