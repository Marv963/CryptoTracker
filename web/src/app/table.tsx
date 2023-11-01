"use client";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Card, CardBody } from "@nextui-org/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { FaAngleDown, FaAngleUp } from "react-icons/fa6";
import { useEffect, useRef, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { motion, useAnimation } from "framer-motion";
import { useTheme } from "next-themes";
import { useRouter } from "next/navigation";

import { Data } from "@/types/symbol";

const placeholderData = Array(100).fill({
  pair: "LOADING...",
  price: 0,
  price1dAgo: 0,
  price1hAgo: 0,
  price7dAgo: 0,
  priceChange1d: 0,
  priceChange1h: 0,
  priceChange7d: 0,
  quote_volume: 0,
});

const PriceChange = ({ value }: { value: number | undefined }) => {
  if (value == undefined) {
    return <></>;
  }
  return (
    <div className="flex justify-end">
      <span
        className={`flex items-center w-[86px] justify-center gap-2 bg-opacity-30 px-2 py-1 rounded-xl  ${
          value < 0
            ? "bg-danger-100 text-danger-400"
            : "bg-success-100 text-success-400"
        } `}
      >
        {value < 0 ? <FaAngleDown /> : <FaAngleUp />}
        {value.toFixed(2)} %
      </span>
    </div>
  );
};

const UpdatedPriceCell = ({
  price,
  initialColor,
}: {
  price: number | undefined;
  initialColor: string;
}) => {
  const prevPriceRef = useRef<number>();
  const controls = useAnimation();

  useEffect(() => {
    if (price == undefined) {
      return;
    }
    controls.set({ color: initialColor });
    // Wenn es keinen vorherigen Preis gibt oder der Preis gleich geblieben ist, nichts tun
    if (prevPriceRef.current === undefined || prevPriceRef.current === price) {
      prevPriceRef.current = price;
      return;
    }

    const comparisonPrice = prevPriceRef.current;
    // Den aktuellen Preis für die nächste Prüfung speichern
    prevPriceRef.current = price;

    // Animation starten
    controls.start({
      color: [
        `${comparisonPrice < price ? "#45d383" : "#f54281"}`,
        // `${comparisonPrice < price ? "#45d383" : "#f54281"}`,
        initialColor,
      ], // Grün zu Schwarz
      transition: { duration: 3 },
    });
  }, [price, controls, initialColor]);
  if (price == undefined) {
    return <></>;
  }

  return (
    <motion.td className="text-right font-extrabold" animate={controls}>
      <span>${price < 10 ? price.toFixed(4) : price.toFixed(2)}</span>
    </motion.td>
  );
};

const MotionCard = motion(Card);

export default function CryptoCurrenciesTable() {
  const [messageSent, setMessageSent] = useState(false);
  const queryClient = useQueryClient();
  const { resolvedTheme } = useTheme();
  const router = useRouter();

  const { data } = useQuery({
    queryKey: ["symbol"],
    queryFn: async () => {
      const { data } = await axios.get("/api/symbols");
      return data as Data[];
    },
  });

  const { sendJsonMessage, readyState } = useWebSocket(
    process.env.NEXT_PUBLIC_WEBSOCKET_URL!,
    {
      onOpen: () => {
        console.log("WebSocket connection established.");
      },
      share: true,
      filter: () => false,
      retryOnError: true,
      shouldReconnect: () => true,
      onMessage: (event) => {
        console.log("Received message: ", event.data);
        handleWebSocketMessage(JSON.parse(event.data));
      },
    },
  );

  useEffect(() => {
    if (readyState === ReadyState.OPEN && !messageSent && data) {
      // Senden Sie Ihre Nachricht hier
      sendJsonMessage({
        method: "subscribe",
        symbols: data.map((d) => d.pair),
      });

      // Setzen Sie Ihren Status auf true, so dass die Nachricht nicht erneut gesendet wird
      setMessageSent(true);
    }
  }, [readyState, messageSent, sendJsonMessage, data]);

  const handleWebSocketMessage = (updatedData: Data) => {
    console.log("get message", updatedData);
    queryClient.setQueryData(["symbol"], (prevData: Data[] | undefined) => {
      if (!prevData) {
        return undefined;
      }

      // Aktualisieren Sie den Datensatz
      return prevData.map((dataItem) =>
        dataItem.pair === updatedData.pair ? updatedData : dataItem,
      );
    });
  };

  const rowsData = data || placeholderData;

  // Definiere deine Animationsvarianten
  const variants = {
    hidden: { opacity: 0, y: 1000 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <MotionCard
      className="w-full bg-opacity-20 shadow-2xl shadow-blue-500 dark:shadow-cyan-400"
      initial="hidden"
      animate="visible"
      variants={variants}
      transition={{ duration: 0.3, delay: 0.5 }} // Optional: Dauer der Animation anpassen
    >
      {/* shadow-red-600*/}
      <CardBody>
        <Table>
          <TableCaption>CryptoCurrencies</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">#</TableHead>
              <TableHead>Name</TableHead>
              <TableHead className="text-right">Price</TableHead>
              <TableHead className="text-right w-[130px]">1h %</TableHead>
              <TableHead className="text-right w-[130px]">24h %</TableHead>
              <TableHead className="text-right w-[130px]">7d %</TableHead>
              {/* <TableHead>Price Graph (7d)</TableHead> */}
            </TableRow>
          </TableHeader>
          <TableBody>
            {rowsData.map((d, i) => (
              <TableRow
                onClick={() =>
                  router.push(`/coins/${d.pair.replace("/", "-")}`)
                }
                key={i}
                className="cursor-pointer"
              >
                <TableCell className="text-muted-foreground">{i + 1}</TableCell>
                <TableCell className="">
                  {d.pair && d.pair.replace("/USD", "")}
                </TableCell>
                <UpdatedPriceCell
                  price={d.price}
                  initialColor={resolvedTheme === "light" ? "#000" : "#fff"}
                />
                <TableCell>
                  <PriceChange value={d.priceChange1h} />
                </TableCell>
                <TableCell>
                  <PriceChange value={d.priceChange1d} />
                </TableCell>
                <TableCell>
                  <PriceChange value={d.priceChange7d} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardBody>
    </MotionCard>
  );
}
