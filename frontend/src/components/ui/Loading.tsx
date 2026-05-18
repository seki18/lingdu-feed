"use client";

import { useEffect, useRef } from "react";
import lottie from "lottie-web";

export default function Loading({ size = 80 }: { size?: number }) {
  const container = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!container.current) return;

    const anim = lottie.loadAnimation({
      container: container.current,
      renderer: "svg",
      loop: true,
      autoplay: true,
      path: "/lottie/loading.json",
    });

    return () => anim.destroy();
  }, []);

  return <div ref={container} style={{ width: size, height: size, margin: "0 auto" }} />;
}
