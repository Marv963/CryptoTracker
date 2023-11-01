from typing import Optional, Type
import ccxt
import pika
from pika.adapters.blocking_connection import BlockingChannel, BlockingConnection
import json
from streamservice.ticker_types import Ticker
from streamservice.config import Config


class MessageQueueClient:
    def __init__(self, config: Type[Config]):
        self.config = config
        self.connection: Optional[BlockingConnection] = None
        self.channel: Optional[BlockingChannel] = None

    def connect(self) -> None:
        credentials = pika.PlainCredentials(
            self.config.RABBITMQ_USERNAME, self.config.RABBITMQ_PASSWORD
        )
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(
                self.config.RABBITMQ_HOST,
                self.config.RABBITMQ_PORT,
                "/",
                credentials,
            )
        )
        self.channel = self.connection.channel()
        self.channel.exchange_declare(
            exchange=self.config.RABBITMQ_EXCHANGE,
            exchange_type="topic",  # oder 'direct', 'fanout', etc. je nach Use-Case
            durable=True,  # oder False, je nachdem, ob der Exchange nach einem Neustart von RabbitMQ erhalten bleiben soll
        )
        self.channel.queue_declare(queue="kraken_queue", durable=True)
        self.channel.queue_bind(
            queue="kraken_queue",
            exchange=self.config.RABBITMQ_EXCHANGE,
            routing_key="kraken_routing_key",
        )

    def write_tickers(self, tickers: Ticker) -> None:
        for _, ticker_data in tickers.items():
            message = json.dumps(ticker_data)
            self.channel.basic_publish(
                exchange=self.config.RABBITMQ_EXCHANGE,
                routing_key="kraken_routing_key",
                body=message,
            )


class CryptoExchange:
    def __init__(self, exchange_name: str):
        self.exchange = getattr(ccxt, exchange_name)()

    def fetch_tickers(self) -> Ticker:
        return self.exchange.fetch_tickers()


Config.validate()

mq = MessageQueueClient(config=Config)
mq.connect()
kraken = CryptoExchange(exchange_name="kraken")
tickers = kraken.fetch_tickers()
mq.write_tickers(tickers=tickers)
