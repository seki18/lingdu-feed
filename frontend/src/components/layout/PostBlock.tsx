import Link from "next/link";

interface Props {
  id: number;
}

export default function PostBlock({ id }: Props) {
  return (
    <Link href={`/posts/${id}`}>
      <div className="cursor-pointer rounded border p-4 shadow-sm transition hover:bg-gray-50">
        <h2 className="mb-2 text-lg font-bold">
          Post {id}
        </h2>

        <p className="text-gray-600">
          Click to open detail page
        </p>
      </div>
    </Link>
  );
}