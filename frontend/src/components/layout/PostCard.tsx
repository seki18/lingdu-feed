import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { PostSummary } from "@/types/post";
import { trackState } from "@/lib/api";

interface Props {
  post: PostSummary;
  showDelimited?: boolean;
}

function Icon({ name }: { name: string }) {
  return <img src={`/icon/${name}.svg`} alt={name} style={{ width: 16, height: 16, display: "inline" }} />;
}

export default function PostCard({ post }: Props) {
  const router = useRouter();

  // Track resource display (status=2) when component mounts
  useEffect(() => {
    trackState(post.id, 2);
  }, [post.id]);

  // Track resource click (status=3) — navigate directly (no URL caching)
  const handleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    trackState(post.id, 3);
    router.push(`/posts/${post.id}`);
  };

  const handleUserClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    router.push(`/users/${post.user_id}`);
  };

  return (
    <div className="mb-4 cursor-pointer rounded border p-4 shadow-sm transition hover:bg-gray-50" onClick={handleClick}>
      <h2 className="text-lg font-bold">{post.title}</h2>
      <div className="mt-2 flex items-center gap-4 text-sm text-gray-500">
        <span
          onClick={handleUserClick}
          className="inline-flex items-center gap-1 text-gray-500 hover:text-blue-600 hover:underline cursor-pointer"
        >
          <Icon name="me" /> {post.username || `User ${post.user_id}`}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="see" /> {post.stats?.view_count ?? 0}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="praise_no" /> {post.stats?.like_count ?? 0}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="comment" /> {post.stats?.comment_count ?? 0}
        </span>
        <span className="inline-flex items-center gap-1">
          <Icon name="collect_no" /> {post.stats?.favorite_count ?? 0}
        </span>
      </div>
      <p className="mt-1 text-sm text-gray-400">{new Date(post.created_time).toLocaleString()}</p>
    </div>
  );
}
