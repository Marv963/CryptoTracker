"use client";

import Heading from "./heading";
import ArbitrageTable from "./table";

export default function Arbitrage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-between gap-10">
      <Heading />
      <ArbitrageTable />
    </div>
  );
}
