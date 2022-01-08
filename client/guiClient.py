from tkinter import *
from tkinter import simpledialog
from guiFuncs import *


class App:
	def __init__(self, root):
		self.root = root

		self.initValues()
		self.refreshWallets()
		self.initWidgets()

	def initValues(self):
		self.wallets = StringVar()
		self.balanceOf = StringVar()
		self.selectedWallet = StringVar()

	def refreshWallets(self):
		self.wallets.set(BlockChain.queryWallets())

	def refreshBalance(self):
		if self.selectedWallet != "":
			self.balanceOf.set("Balance: {}".format(BlockChain.getBalanceOf(self.selectedWallet.get())))


	def initWidgets(self):
		# initWalletsDiv
		self.initWalletsDiv()

		# initTransactionDiv
		self.initTransactionDiv()

	# initWalletsDiv : a div of wallet
	def initWalletsDiv(self):
		# walletsDiv && scroll
		walletsDiv = Frame(self.root)
		walletsDiv.pack(side=LEFT, fill=BOTH, expand=YES)

		## label && wallets:listBox && scroll && button
		Label(walletsDiv, text='Wallets', font=24).pack(side=TOP)
		lb = Listbox(walletsDiv, font=24, listvariable=self.wallets)
		scroll = Scrollbar(walletsDiv, command=lb.yview)
		lb.configure(yscrollcommand=scroll.set); scroll.pack(side=RIGHT, fill=Y)
		lb.bind("<<ListboxSelect>>", self.selectWallet(lb))
		lb.pack(side=TOP, padx=(10,0))
		Button(walletsDiv, text='add wallet', command=self.addWallet).pack(side=TOP)

	# initTransactionDiv : a div of transaction
	def initTransactionDiv(self):
		# transactionDiv
		transactionDiv = Frame(self.root)
		transactionDiv.pack(side=LEFT, fill=BOTH, expand=YES)
		# title(label) && selectTip(label) && tip/input && confirm button
		Label(transactionDiv, text='Transaction', font=24).pack(side=TOP)
		self.tipDiv = Label(transactionDiv, text='Please select one wallet', font=24)  # need to be removed by operations

		self.makeTransactionDiv = Frame(transactionDiv)
		self.balanceTipDiv = Label(self.makeTransactionDiv, font=('StSong', 14), textvariable=self.balanceOf).pack(side=TOP)

		# current wallet
		self.transactionCurWalletDiv = Label(self.makeTransactionDiv, textvariable=self.selectedWallet).pack(side=TOP)

		# to wallet
		self.transactionWalletTipDiv = Label(self.makeTransactionDiv, text="to: ").pack(side=LEFT)
		self.transactionWalletValDiv = Entry(self.makeTransactionDiv)
		self.transactionWalletValDiv.pack(side=LEFT, padx=(0, 10))
		# value
		self.transactionValTipDiv = Label(self.makeTransactionDiv, text="val: ").pack(side=LEFT)
		self.transactionValValDiv = Entry(self.makeTransactionDiv)
		self.transactionValValDiv.pack(side=LEFT, padx=(0, 10))
		# miner wallet
		self.transactionMinerTipDiv = Label(self.makeTransactionDiv, text="miner: ").pack(side=LEFT)
		self.transactionMinerValDiv = Entry(self.makeTransactionDiv)
		self.transactionMinerValDiv.pack(side=LEFT, padx=(0, 10))

		# confirm
		self.transactionButtonDiv = Button(transactionDiv, text="make transaction", command=self.makeTransaction(transactionDiv))
		if self.selectedWallet.get() == "":
			self.tipDiv.pack(side=TOP)

		# init block-chain button
		self.initBlockChainButtonDiv = Button(transactionDiv, text="init blockChain", command=self.initBlockChain).pack(side=BOTTOM)

	# define envent
	# when changing selected wallet
	def selectWallet(self, lb):
		def inner(event):
			if len(lb.curselection()) == 1:
				self.selectedWallet.set(lb.get([lb.curselection()[0]]))
				self.tipDiv.pack_forget()
				# transaction info
				self.makeTransactionDiv.pack()
				self.transactionButtonDiv.pack(side=TOP)
				self.refreshBalance()
		return inner

	# when tap addWallet button
	def addWallet(self):
		BlockChain.addWallet()
		self.refreshWallets()

	# when tap make Transaction button
	def makeTransaction(self, transavtionDiv):
		def inner():
			from_ = self.selectedWallet.get()
			to_ = self.transactionWalletValDiv.get()
			val_ = self.transactionValValDiv.get()
			miner_ = self.transactionMinerValDiv.get()

			res = BlockChain.makeTransaction(from_, to_, val_, miner_)
			simpledialog.SimpleDialog(transavtionDiv, title='success', text=res, buttons=["sure"], default=0).go()
			self.refreshBalance()
		return inner

	# when tap init block chain button
	def initBlockChain(self):
		BlockChain.initBlockChain(self.selectedWallet.get())
		self.refreshBalance()


if __name__ == "__main__":
	root = Tk()
	root.title("blockChain-By Cookie")
	display = App(root)
	root.mainloop()
