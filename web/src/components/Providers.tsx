"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "next-themes";
import { NextUIProvider } from "@nextui-org/react";

const Providers = ({ children }: { children: React.ReactNode }) => {
  const queryClient = new QueryClient();
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <NextUIProvider>
        <QueryClientProvider client={queryClient}>
          {children}
        </QueryClientProvider>
      </NextUIProvider>
    </ThemeProvider>
  );
};

export default Providers;
