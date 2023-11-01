import { useEffect } from "react";
// import styled from "styled-components";
import { useAnimation, motion } from "framer-motion";
import { useInView } from "react-intersection-observer";

interface Props {
  title: string;
}

export default function AnimatedTitle({ title }: Props) {
  const ctrls = useAnimation();

  const { ref, inView } = useInView({
    threshold: 0.5,
    triggerOnce: true,
  });

  useEffect(() => {
    if (inView) {
      ctrls.start("visible");
    }
    if (!inView) {
      ctrls.start("hidden");
    }
  }, [ctrls, inView]);

  const wordAnimation = {
    hidden: {},
    visible: {},
  };

  const characterAnimation = {
    hidden: {
      opacity: 0,
      y: `0.25em`,
    },
    visible: {
      opacity: 1,
      y: `0em`,
      transition: {
        duration: 1,
        ease: [0.2, 0.65, 0.3, 0.9],
      },
    },
  };

  return (
    <h2 aria-label={title} role="heading" className="text-5xl">
      {title.split(" ").map((word, index) => (
        <motion.span
          ref={ref}
          aria-hidden="true"
          key={index}
          initial="hidden"
          animate={ctrls}
          variants={wordAnimation}
          className="mr-3"
          transition={{
            delayChildren: index * 0.25,
            staggerChildren: 0.05,
          }}
        >
          {word.split("").map((character, index) => (
            <motion.span
              aria-hidden="true"
              key={index}
              variants={characterAnimation}
            >
              {character}
            </motion.span>
          ))}
        </motion.span>
      ))}
    </h2>
  );
}
