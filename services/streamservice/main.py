import argparse
from typing import Optional, Type
import ccxt
import pika
from pika.adapters.blocking_connection import BlockingChannel, BlockingConnection
import json
import time

from .ticker_types import Ticker
from .config import Config
from .logger import Logger


class CryptoExchange:
    def __init__(self, exchange_name: str):
        exchange_class = getattr(ccxt, exchange_name)
        self.exchange = exchange_class(
            {
                "enableRateLimit": True,
            }
        )

    def fetch_tickers(self) -> Ticker:
        return self.exchange.fetch_tickers()


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

        # Exchange deklarieren
        self.channel.exchange_declare(
            exchange=self.config.RABBITMQ_EXCHANGE,
            exchange_type="topic",  # oder 'direct', 'fanout', etc. je nach Use-Case
            durable=True,  # oder False, je nachdem, ob der Exchange nach einem Neustart von RabbitMQ erhalten bleiben soll
        )

        self.channel.queue_declare(queue=self.config.RABBITMQ_QUEUE, durable=True)
        self.channel.queue_bind(
            queue=self.config.RABBITMQ_QUEUE,
            exchange=self.config.RABBITMQ_EXCHANGE,
            routing_key=self.config.RABBITMQ_ROUTING_KEY,
        )

    def close_connection(self) -> None:
        if self.connection:
            self.connection.close()

    def write_tickers(self, tickers: Ticker) -> None:
        try:
            # Hier sicherstellen, dass self.channel nicht None ist
            assert self.channel is not None, "Channel is None"

            for ticker_name, ticker_data in tickers.items():
                # Convert the ticker data to a string, e.g., via JSON serialization
                message = json.dumps(ticker_data)
                self.channel.basic_publish(
                    exchange=self.config.RABBITMQ_EXCHANGE,
                    routing_key=self.config.RABBITMQ_ROUTING_KEY,
                    body=message,
                    # properties=pika.BasicProperties(delivery_mode=2),
                )
                Logger.debug(f"Sent ticker data for {ticker_name} to RabbitMQ")
        except Exception as e:
            Logger.error(f"Failed to send message to RabbitMQ: {str(e)}")


def main(exchange_name: str) -> None:
    Logger.initialize_logger("streamservice.log")
    Logger.info("Start")
    try:
        Config.validate()
    except ValueError as e:
        Logger.error(f"Configuration error: {str(e)}")
        exit(1)

    # Der Producer sendet die Daten an einen TopicExchange mit den Routing Schlüssel des Börsennamens
    Config.RABBITMQ_QUEUE = f"queue_{exchange_name}_prices"
    Config.RABBITMQ_ROUTING_KEY = f"{exchange_name}.prices"

    data_writer = MessageQueueClient(config=Config)
    try:
        # Überprüfen der Verbindung zu RabbitMQ
        data_writer.connect()
        Logger.info("Connected to RabbitMQ")
    except Exception as e:
        Logger.error(f"Failed to connect to RabbitMQ: {str(e)}")
        exit(1)

    # Instanz von CryptoExchange für das ausgewählte Exchange
    exchange = CryptoExchange(exchange_name)

    print("Fetching ticker data from {}... Press Ctrl+C to stop.".format(exchange_name))

    try:
        while True:
            try:
                tickers = exchange.fetch_tickers()
                data_writer.write_tickers(tickers)
                time.sleep(5)
            except Exception as e:
                Logger.error(f"An error occurred: {str(e)}")
                data_writer.close_connection()
                time.sleep(60)
                data_writer.connect()
                time.sleep(20)
    except KeyboardInterrupt:
        print("\nKeyboard interrupt received. Stopping...")

    finally:
        # Sicherstellen, dass die Verbindung ordnungsgemäß geschlossen wird, auch wenn ein Fehler auftritt
        data_writer.close_connection()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Fetch tickers from a specified cryptocurrency exchange."
    )
    parser.add_argument(
        "exchange", type=str, help="Name of the cryptocurrency exchange"
    )

    args = parser.parse_args()
    main(args.exchange)
