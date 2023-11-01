"use client";
import AnimatedTitle from "@/components/AnimatedTitle";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { AreaData } from "lightweight-charts";
import CoinChart from "./chart";

interface Data {
  exchange: string;
  data: AreaData[];
}

interface Props {
  params: { coin: string };
}

export default function Coin({ params }: Props) {
  const { data, isLoading } = useQuery({
    queryKey: ["symbol"],
    queryFn: async () => {
      const { data } = await axios.get(`/api/symbolhistory/${params.coin}`);
      return data as Data[];
    },
  });
  console.log(data);

  return (
    <div className="flex flex-col gap-12">
      <div className="text-center">
        <AnimatedTitle title={params.coin.replace("-", "/")} />
      </div>
      <div className="flex flex-row">
        <div className="basis-2/3">
          {isLoading ? <div>Isloading...</div> : <CoinChart data={data!} />}
        </div>
      </div>
      <div>
        {/* <Card> */}
        {/*   <CardBody> */}
        {/*     {/*TODO: Hier die preisdifferenzen der letzten tage anzeigen*/}
        {/*   </CardBody> */}
        {/* </Card> */}
      </div>
    </div>
  );
}
