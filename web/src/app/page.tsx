"use client";
import Heading from "./heading";
import CryptoCurrenciesTable from "./table";

export default function CryptoCurrencies() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-between gap-10">
      <Heading />
      <CryptoCurrenciesTable />
    </div>
  );
}
