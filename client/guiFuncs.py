import os
import re
import subprocess
import sys
from typing import *

current_path = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
exe_file = "main.exe" if sys.platform == "win32" else "main"
exe_path = os.path.join(current_path, exe_file)


def runCMD(paramsList: List) -> str:
    cmdList = [exe_path] + paramsList
    r = subprocess.Popen(cmdList, stdout=subprocess.PIPE).communicate()[0]  # ignore_security_alert
    return r.decode().replace("\n", "")


class BlockChain:
    @staticmethod
    def queryWallets() -> tuple:
        res = tuple(runCMD(["la"]).strip("[]").split(" "))
        return res

    @staticmethod
    def addWallet():
        runCMD(["aw"])

    @staticmethod
    def getBalanceOf(address: str) -> float:
        res = runCMD(["gbo", address])
        valList = re.findall(r"\d+\.?\d*$", res)
        if len(valList) == 0:
            return 0
        return float(valList[0])

    @staticmethod
    def makeTransaction(from_: str, to_: str, val_: str, miner_: str) -> str:
        a = runCMD(["mt", from_, to_, val_, miner_])
        return a.replace("【", "\n【").rstrip("=")


    @staticmethod
    def initBlockChain(godAddress: str):
        runCMD(["initbc", godAddress])
