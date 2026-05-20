import { PostSummary } from "./post";

export interface User {
  id: number;
  username: string;
  email: string;
  following_count: number;
  follower_count: number;
  is_following: boolean;
}

export interface UserProfilePage {
  user: User;
  posts: PostSummary[];
  total: number;
  page: number;
  page_size: number;
}