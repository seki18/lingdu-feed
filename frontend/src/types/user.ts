export interface User {
  id: number;
  username: string;
  email: string;
  following_count: number;
  follower_count: number;
  is_following: boolean;
}