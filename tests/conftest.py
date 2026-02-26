import os
import grpc
import pytest
import mysql.connector

@pytest.fixture(scope="session")
def order_channel():
    addr = os.getenv("ORDER_SERVICE_ADDR", "localhost:50052")
    channel = grpc.insecure_channel(addr)
    yield channel
    channel.close()

@pytest.fixture(scope="session")
def mysql_conn():
    conn = mysql.connector.connect(
        host=os.getenv("MYSQL_HOST", "127.0.0.1"),
        port=int(os.getenv("MYSQL_PORT", "3306")),
        user=os.getenv("MYSQL_USER", "minigate"),
        password=os.getenv("MYSQL_PASSWORD", "minigate"),
        database=os.getenv("MYSQL_DB", "minigate"),
    )
    yield conn
    conn.close()
