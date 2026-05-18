export interface CommentItem {
  id: number;
  post_id: number;
  user_id: number;
  username: string;
  reply_id: number | null;
  reply_username: string | null;
  content: string;
  created_time: string;
}

export interface CreateCommentRequest {
  post_id: number;
  content: string;
  reply_id?: number | null;
}
