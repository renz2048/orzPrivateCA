import argparse

class caArgParser():
    def __init__(self) -> None:
        self.opt = None
        self.result = 0

    def parseArgs(self) -> None:
        parser = argparse.ArgumentParser(
            description='用于测试的私有CA工具',
            formatter_class=argparse.RawDescriptionHelpFormatter
        )

        subparsers = parser.add_subparsers(help='证书用途')

        ca_parser = subparsers.add_parser(
            'ca',
            description='创建新的二级CA, 包含orz root CA和orz sub CA',
            help='创建新的二级CA',
            formatter_class=argparse.RawDescriptionHelpFormatter,
            epilog='''
示例:
    python3 %(prog)s
''')
        ca_parser.set_defaults(func=self._ca_handler)

        user_parser = subparsers.add_parser(
            'user',
            description='创建用户证书',
            help='创建用户证书',
            formatter_class=argparse.RawDescriptionHelpFormatter,
            epilog='''
示例:
    python3 %(prog)s
''')
        user_parser.add_argument(
            'common_name',
            help='用户证书CN项(common name)'
        )
        user_parser.set_defaults(func=self._user_handler)

        self.opt = parser.parse_args()

    def _ca_handler(self) -> None:
        pass
    def _user_handler(self) -> None:
        pass