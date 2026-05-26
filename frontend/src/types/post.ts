import { CommentItem } from './comment';

export interface PostStats {
  like_count: number;
  comment_count: number;
  favorite_count: number;
  view_count: number;
}

export interface PostSummary {
  id: number;
  user_id: number;
  username: string;
  title: string;
  created_time: string;
  stats?: PostStats;
  has_liked?: boolean;
  has_favorited?: boolean;
}

export interface Post {
  id: number;
  user_id: number;
  username: string;
  title: string;
  content: string;
  created_time: string;
  updated_time: string;
  stats?: PostStats;
}

export interface PostDetailResponse {
  post: Post;
  has_liked: boolean;
  has_favorited: boolean;
  comments: CommentItem[];
}
