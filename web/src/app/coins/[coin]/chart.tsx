"use client";
import { Card, Select, SelectItem, Selection } from "@nextui-org/react";
import { motion } from "framer-motion";
import {
  createChart,
  ColorType,
  AreaData,
  CrosshairMode,
} from "lightweight-charts";
import React, { useEffect, useRef, useState } from "react";

interface Data {
  exchange: string;
  data: AreaData[];
}

interface ColorScheme {
  backgroundColor?: string;
  lineColor?: string;
  textColor?: string;
  areaTopColor?: string;
  areaBottomColor?: string;
}

interface CoinChartProps {
  data: Data[];
  colors?: ColorScheme;
}

const colors = ["#E63946", "#1D3557", "#2A9D8F"];

const MotionCard = motion(Card);

const CoinChart: React.FC<CoinChartProps> = ({
  data,
  colors: {
    backgroundColor = "white",
    lineColor = "#2962FF",
    textColor = "black",
    areaTopColor = "#2962FF",
    areaBottomColor = "rgba(41, 98, 255, 0.28)",
  } = {},
}) => {
  const [visibleExchanges, setVisibleExchanges] = useState(
    new Set(["average"]),
  );

  const chartContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const { current: chartContainer } = chartContainerRef;
    if (!chartContainer) {
      return;
    }

    const handleResize = () => {
      chart.applyOptions({ width: chartContainer.clientWidth });
    };

    const chart = createChart(chartContainer, {
      width: chartContainer.clientWidth,
      height: 300,
      layout: {
        background: { type: ColorType.Solid, color: backgroundColor },
        textColor,
      },
      grid: {
        horzLines: {
          color: "#F0F3FA",
        },
        vertLines: {
          color: "#F0F3FA",
        },
      },
      crosshair: {
        mode: CrosshairMode.Normal,
      },
      timeScale: {
        borderColor: "rgba(197, 203, 206, 1)",
      },
      handleScroll: {
        vertTouchDrag: false,
      },
    });

    data[0].data &&
      chart.timeScale().setVisibleLogicalRange({
        from: data[0].data.length - 100,
        to: data[0].data.length + 14,
      });

    data.forEach((exchangeData, idx) => {
      if (!visibleExchanges.has(exchangeData.exchange)) {
        return;
      }

      const newSeries = chart.addLineSeries({
        lineWidth: 2,
        color: colors[idx],
      });
      newSeries.setData(exchangeData.data);
    });

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);

      chart.remove();
    };
  }, [
    data,
    backgroundColor,
    lineColor,
    textColor,
    areaTopColor,
    areaBottomColor,
    visibleExchanges,
  ]);

  // Definiere deine Animationsvarianten
  const variants = {
    hidden: { opacity: 0, y: 1000 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <MotionCard
      className="w-full bg-opacity-50 shadow-2xl p-4"
      // className="w-full bg-opacity-20 shadow-2xl shadow-blue-500 dark:shadow-cyan-400 p-4"
      initial="hidden"
      animate="visible"
      variants={variants}
      transition={{ duration: 0.3, delay: 0.5 }} // Optional: Dauer der Animation anpassen
    >
      <div className="flex flex-col gap-4">
        <div>
          <Select
            label="Exchanges"
            selectionMode="multiple"
            placeholder="Select an animal"
            selectedKeys={visibleExchanges}
            className="max-w-xs"
            onSelectionChange={(selection: Selection) =>
              setVisibleExchanges(selection as Set<string>)
            }
          >
            {data.map((d) => (
              <SelectItem key={d.exchange} value={d.exchange}>
                {d.exchange}
              </SelectItem>
            ))}
          </Select>
        </div>
        <div ref={chartContainerRef} className="w-full" />
        <div className="flex gap-4">
          {data.map(
            (d, idx) =>
              visibleExchanges.has(d.exchange) && (
                <div key={d.exchange} className="flex items-center">
                  <div
                    style={{ backgroundColor: colors[idx] }}
                    className="w-4 h-4 mr-2 rounded-full"
                  ></div>
                  <span>{d.exchange}</span>
                </div>
              ),
          )}
        </div>
      </div>
    </MotionCard>
  );
};

export default CoinChart;
