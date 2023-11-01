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
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Input,
  InputProps,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { motion, useAnimation } from "framer-motion";
import { useTheme } from "next-themes";
import { useRouter } from "next/navigation";

import { Data } from "@/types/symbol";
import { calculateProfit } from "./fees";
import { Scale } from "lucide-react";

const InputField: React.FC<InputProps> = (props) => (
  <Input
    isClearable
    radius="lg"
    classNames={{
      label: "text-white/90",
      input: ["bg-transparent", "text-white/90", "placeholder:text-white/60"],
      innerWrapper: "bg-transparent",
      inputWrapper: [
        "shadow-xl",
        "bg-default-200/50",
        "dark:bg-default/60",
        "backdrop-blur-xl",
        "backdrop-saturate-200",
        "hover:bg-default-200/70",
        "focus-within:!bg-default-200/50",
        "dark:hover:bg-default/70",
        "dark:focus-within:!bg-default/60",
        "!cursor-text",
      ],
    }}
    // className="min-w-[200px]"
    type="number"
    {...props}
  />
);

const PriceCell = ({ value }: { value: number | undefined }) => {
  if (value == undefined) {
    return <></>;
  }
  return (
    <div className="text-right">
      {value < 10 ? value.toFixed(4) : value.toFixed(2)}
    </div>
  );
};

const UpdatedPriceCell = ({
  price,
  initialColor,
}: {
  price: number;
  initialColor: string;
}) => {
  const prevPriceRef = useRef<number>();
  const controls = useAnimation();

  useEffect(() => {
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
        initialColor,
      ], // Grün zu Schwarz
      transition: { duration: 3 },
    });
  }, [price, controls, initialColor]);

  return (
    <motion.td className="text-right font-extrabold" animate={controls}>
      <span>${price < 10 ? price.toFixed(4) : price.toFixed(2)}</span>
    </motion.td>
  );
};

const MotionCard = motion(Card);

export default function ArbitrageTable() {
  const [messageSent, setMessageSent] = useState(false);
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const queryClient = useQueryClient();
  const { theme } = useTheme();
  const router = useRouter();
  const [amount, setAmount] = useState(0);
  const [tradeVolumeLowerExchange, setTradeVolumeLowerExchange] = useState(0);
  const [tradeVolumeHigherExchange, setTradeVolumeHigherExchange] = useState(0);
  const [selectedData, setSelectedData] = useState<Data | null>(null);

  const { data } = useQuery({
    queryKey: ["symbol"],
    queryFn: async () => {
      const { data } = await axios.get("/api/arbitrage");
      return data as Data[];
    },
  });

  const { sendJsonMessage, readyState } = useWebSocket(
    "wss://crypto-tracker.app/ws",
    {
      onOpen: () => {
        console.log("WebSocket connection established.");
      },
      share: true,
      filter: () => false,
      retryOnError: true,
      shouldReconnect: () => true,
      onMessage: (event) => {
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

  const placeholderData = Array(100).fill({
    pair: "LOADING...",
    price: 0,
    lowestPrice: 0,
    highestPrice: 0,
    arbitrage: 0,
  });

  const rowsData = data || placeholderData;

  const handleRowClick = (data: Data) => {
    setSelectedData(data);
    onOpen(); // Öffne das Modal hier, nachdem die Daten gesetzt wurden
  };

  const handleAmountChange = (value: number) => {
    // Hier können Sie auch Validierungen hinzufügen
    setAmount(value);
  };

  // const calculateProfit = (amount: number) => {
  //   return 0;
  //   // Ihre Logik zur Berechnung des Gewinns...
  //   // Beispiel: return amount * currentArbitragePercentage;
  // };

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
              <TableHead className="text-right w-[130px]">
                Average Price
              </TableHead>
              <TableHead className="text-right w-[130px]">
                Lowest Price
              </TableHead>
              <TableHead className="text-right w-[130px]">
                Highest Price
              </TableHead>
              <TableHead className="text-right w-[130px]">
                Differenz in %
              </TableHead>
              <TableHead className="text-right w-[170px]">
                Calculated Arbitrage
              </TableHead>
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
                <TableCell className="">{d.pair}</TableCell>
                {d.price ? (
                  <UpdatedPriceCell
                    price={d.price}
                    initialColor={theme === "dark" ? "#fff" : "#000"}
                  />
                ) : (
                  <></>
                )}
                <TableCell>
                  <PriceCell value={d.lowestPrice} />
                </TableCell>
                <TableCell>
                  <PriceCell value={d.highestPrice} />
                </TableCell>
                <TableCell>
                  <PriceCell value={d.arbitrage} />
                </TableCell>
                <TableCell className="flex justify-center">
                  <Button
                    radius="full"
                    className="bg-gradient-to-tr from-blue-500 to-cyan-400 text-white shadow-lg"
                    onPress={() => handleRowClick(d)}
                  >
                    Caluclate
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
          <Modal isOpen={isOpen} onOpenChange={onOpenChange} backdrop="blur">
            <ModalContent className="dark bg-gradient-to-tr from-blue-500 to-cyan-400 text-white/90 w-[420px]">
              {(onClose) => (
                <>
                  <ModalHeader className="flex flex-col gap-1">
                    Calculate Arbitrage Trade
                  </ModalHeader>
                  <ModalBody className="">
                    <Card className="bg-transparent shadow-md">
                      <CardHeader>{selectedData?.pair}</CardHeader>
                      <CardBody className="flex gap-4">
                        {/* Price Information */}
                        <div className="flex justify-between space-x-4">
                          {/* Container für beide Karten */}
                          {/* Karte für Lowest Price */}
                          <div className="flex flex-col px-4 py-2 rounded-xl bg-blue-500">
                            <p className="text-sm text-white/80">
                              Lowest Price
                            </p>
                            {/* Kleiner Text */}
                            <p className="text-2xl">
                              {selectedData?.lowestPrice &&
                              selectedData?.lowestPrice < 10
                                ? selectedData.lowestPrice.toFixed(4)
                                : selectedData?.lowestPrice.toFixed(2)}
                            </p>
                            {/* Großer Text */}
                            <p className="text-sm text-white/80">
                              at {selectedData?.lowestPriceExchange}
                            </p>
                            {/* Kleiner Text */}
                          </div>
                          <div className="flex flex-col px-4 py-2 justify-center align-middle">
                            {/* <div className="w-20 h-20 bg-sky-500"></div> */}
                            <p>
                              {selectedData &&
                              selectedData.arbitrage &&
                              selectedData.arbitrage < 0
                                ? ""
                                : "+"}
                              {selectedData!.arbitrage.toFixed(0)}%
                            </p>
                          </div>
                          {/* Karte für Highest Price */}
                          <div className="flex flex-col px-4 py-2 rounded-xl bg-blue-500">
                            <p className="text-sm text-white/80">
                              Highest Price
                            </p>
                            <p className="text-2xl">
                              {selectedData?.highestPrice &&
                              selectedData?.highestPrice < 10
                                ? selectedData?.highestPrice.toFixed(4)
                                : selectedData?.highestPrice.toFixed(2)}
                            </p>
                            <p className="text-sm text-white/80">
                              at {selectedData?.highestPriceExchange}
                            </p>
                          </div>
                        </div>
                        <InputField
                          label={`Amount in ${selectedData?.pair.split(
                            "/",
                          )[1]}`}
                          placeholder={`Amount of ${selectedData?.pair.split(
                            "/",
                          )[1]} to Trade`}
                          onChange={(e) => setAmount(+e.target.value)}
                          startContent={
                            <Scale className="text-white/90 text-slate-400 pointer-events-none flex-shrink-0" />
                          }
                        />

                        <InputField
                          label={`Trade Volume on ${selectedData?.lowestPriceExchange} last 30 Days`}
                          placeholder={`Trade Volume on ${selectedData?.lowestPriceExchange} last 30 Days`}
                          onChange={(e) =>
                            setTradeVolumeLowerExchange(+e.target.value)
                          }
                        />
                        <InputField
                          label={`Trade Volume on ${selectedData?.highestPriceExchange} last 30 Days`}
                          placeholder={`Trade Volume on ${selectedData?.highestPriceExchange} last 30 Days`}
                          onChange={(e) =>
                            setTradeVolumeHigherExchange(+e.target.value)
                          }
                        />
                        {/* Fee and Profit */}
                        <div className="text-2xl">
                          Estimated Profit: {selectedData?.pair.split("/")[1]}{" "}
                          {calculateProfit(
                            selectedData?.pair,
                            selectedData?.lowestPrice,
                            selectedData?.highestPrice,
                            amount,
                            tradeVolumeLowerExchange,
                            tradeVolumeHigherExchange,
                            selectedData?.lowestPriceExchange,
                            selectedData?.highestPriceExchange,
                          )}{" "}
                          {/* <p>Estimated Profit: ${calculateProfit(amount)} </p> */}
                        </div>

                        {/* Arbitrage Information */}
                        <div className="flex flex-col gap-2">
                          {/* <p> */}
                          {/*   Arbitrage Opportunity: ${selectedData?.arbitrage} */}
                          {/* </p> */}
                          <p>
                            <span className="font-bold">Step 1:</span> Buy{" "}
                            {selectedData?.pair} at{" "}
                            {selectedData?.lowestPriceExchange} for{" "}
                            {selectedData?.lowestPrice}.
                          </p>
                          <p>
                            <span className="font-bold">Step 2:</span> Sell at{" "}
                            {selectedData?.highestPriceExchange}
                            for ${selectedData?.highestPrice}.
                          </p>
                        </div>
                      </CardBody>
                    </Card>
                  </ModalBody>
                  <ModalFooter>
                    <Button color="danger" variant="light" onPress={onClose}>
                      Close
                    </Button>
                  </ModalFooter>
                </>
              )}
            </ModalContent>
          </Modal>
        </Table>
      </CardBody>
    </MotionCard>
  );
}
