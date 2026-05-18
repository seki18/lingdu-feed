export interface PostSummary {
  id: number;
  user_id: number;
  username: string;
  title: string;
  created_time: string;
  praise_count: number;
  comment_count: number;
  collection_count: number;
  view_count: number;
  has_praised?: boolean;
  has_collected?: boolean;
}

export interface PostDetail {
  id: number;
  user_id: number;
  username: string;
  title: string;
  content: string;
  created_time: string;
  updated_time: string;
}
