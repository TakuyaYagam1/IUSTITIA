import argparse

from api.http.handlers.handlers import Server

if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("--host", default="0.0.0.0")
    parser.add_argument("--port", default=9000)

    parser.add_argument("--operations_host", type=str)
    parser.add_argument("--operations_port", type=int)
    parser.add_argument(
        "--standalone",
        action="store_true",
        help="skip /game/checker/register handshake (for local dev).",
    )

    args = parser.parse_args()

    operations_hostport = ""
    if args.operations_host and args.operations_port:
        operations_hostport = args.operations_host + ":" + str(args.operations_port)

    server = Server(
        host=str(args.host),
        port=int(args.port),
        operations_hostport=operations_hostport,
        standalone=bool(args.standalone),
    )
    server.run()
