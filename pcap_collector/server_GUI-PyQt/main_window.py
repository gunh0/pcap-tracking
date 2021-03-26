import os
import sys
import time

from queue import Queue
from threading import Thread
from command import command
from command import make_moz_cmd_list, make_tbb_cmd_list, make_copy_cmd_list, cmd_cmd_list, exit_cmd_list, watch_cmd_list

import file_open_easygui as fopen

from PyQt5 import QtCore, QtGui, QtWidgets
from PyQt5 import uic

from tkinter import Tk, Listbox, Label

# import .ui file
# .ui file must be located in the same directory as the Python code file.
form_class = uic.loadUiType("CC_Server_GUI.ui")[0]

print("Running in: ", os.getcwd())


class CommandSerize9999(Thread):
    def __init__(self, queue, ip_addr):
        Thread.__init__(self)
        self.queue = queue
        self.ip_addr = ip_addr
        self.port = "9999"

    def run(self):
        while True:
            command_list = self.queue.get()
            if command_list is None:
                break
            for c in range(command_list.qsize()):
                com = command_list.get()
                print(c, ": ", com)
                command(self.ip_addr, self.port, com)
                time.sleep(0.5)
            self.queue.task_done()


class CommandSerize9998(Thread):
    def __init__(self, queue, ip_addr):
        Thread.__init__(self)
        self.queue = queue
        self.ip_addr = ip_addr
        self.port = "9998"

    def run(self):
        while True:
            command_list = self.queue.get()
            if command_list is None:
                break

            for c in range(command_list.qsize()):
                com = command_list.get()
                print(c, ": ", com)
                command(self.ip_addr, self.port, com)

                time.sleep(0.5)
            self.queue.task_done()


class WindowClass(QtWidgets.QMainWindow, form_class):   # GUI Class Define
    def __init__(self):
        super().__init__()
        self.setupUi(self)

        # 100 Tab Func.
        self.pBtnSet100ip.clicked.connect(self.pBtnSet100ip_function)
        self.pBtn100_collect.clicked.connect(self.pBtn100_collect_function)
        self.pBtn100_cmd.clicked.connect(self.pBtn100_cmd_function)
        self.pBtn100_copy.clicked.connect(self.pBtn100_copy_function)

        # 200 Tab Func.
        self.pBtnSet200ip.clicked.connect(self.pBtnSet200ip_function)
        self.pBtn200_watch.setEnabled(False)
        self.pBtn200_copy.clicked.connect(self.pBtn200_copy_function)

        # GCP Tab Func.
        self.pBtnSetGCPip.clicked.connect(self.pBtnSetGCPip_function)
        self.pBtnGCP_OpenTargetList.clicked.connect(
            self.pBtnGCP_OpenTargetList_function)   # Open Tor URL List Text
        self.pBtnGCP_cmd.clicked.connect(self.pBtnGCP_cmd_function)
        self.pBtnGCP_collect.clicked.connect(self.pBtnGCP_collect_function)
        self.pBtnGCP_watch.clicked.connect(self.pBtnGCP_watch_function)
        self.pBtnGCP_copy.clicked.connect(self.pBtnGCP_copy_function)
        self.pBtnGCP_send.setEnabled(False)
        self.pBtnGCP_close.setEnabled(False)    # for close?

    # Share
    ip_list = []
    cmd_list = []
    thread_list = []
    taskQueue = Queue()

    # 100 Tab Function
    def pBtnSet100ip_function(self):
        print("100 IP List Set Buttun Pressed.")
        self.ip_list = []   # init. list
        try:
            start_ip_host = int(self.sip1004.text())
            last_ip_host = int(self.lip1004.text())
            for i in range(start_ip_host, last_ip_host+1):
                input_ip = self.sip1001.text()+"."+self.sip1002.text()+"." + \
                    self.sip1003.text()+"."+str(i)
                self.ip_list.append(input_ip)
            print("IP List: ", self.ip_list)
        except:
            print("Input IP.")

        tkWindow = Tk()     # show IP List with Tk
        tkWindow.geometry("220x200")
        tkWindow.title("IP LIST")
        labeltxt = "Number of Collectors: "+str(len(self.ip_list))
        tkLabel = Label(tkWindow, text=labeltxt)
        tkLabel.pack()
        ipListBox = Listbox(tkWindow)
        for i in range(0, len(self.ip_list)):
            ipListBox.insert(i+1, self.ip_list[i])
        ipListBox.pack()
        tkWindow.mainloop()
        print("100 IP List Set Buttun Finish.")

    def pBtn100_collect_function(self):     # NOT IMPLEMENTED!!!
        print("100 Collect Pressed")
        print("100 Collect Finish")

    def pBtn100_cmd_function(self):
        print("100 CMD Sending Pressed")

        send_command = self.cmd100.text()
        self.cmd_list = []  # init. Command List
        self.taskQueue = Queue()  # init. Task Queue
        for i in self.ip_list:
            self.cmd_list.append(cmd_cmd_list(send_command))

        for item in self.cmd_list:
            self.taskQueue.put(item)

        for ip_addr in self.ip_list:
            t = CommandSerize9998(self.taskQueue, ip_addr)
            t.start()
        self.taskQueue.join()
        print("Queue: ", self.taskQueue)

        for i in range(len(self.ip_list)):
            self.taskQueue.put(None)

        for t in self.thread_list:
            t.join()

        print("100 CMD Sending Finish")

    def pBtn100_copy_function(self):
        print("100 Copy Pressed")
        self.cmd_list = []  # init. Command List
        self.taskQueue = Queue()  # init. Task Queue
        for i in self.ip_list:
            self.cmd_list.append(make_copy_cmd_list())

        for item in self.cmd_list:
            self.taskQueue.put(item)

        for ip_addr in self.ip_list:
            t = CommandSerize9999(self.taskQueue, ip_addr)
            t.start()
        self.taskQueue.join()
        print("Queue: ", self.taskQueue)

        for i in range(len(self.ip_list)):
            self.taskQueue.put(None)

        for t in self.thread_list:
            t.join()

        print("100 Copy Finish")

    # 200 Tab Function
    def pBtnSet200ip_function(self):
        print("200 IP List Set")
        self.ip_list = []   # init. list
        try:
            start_ip_host = int(self.sip2004.text())
            last_ip_host = int(self.lip2004.text())
            for i in range(start_ip_host, last_ip_host+1):
                input_ip = self.sip2001.text()+"."+self.sip2002.text()+"." + \
                    self.sip2003.text()+"."+str(i)
                self.ip_list.append(input_ip)
            print("IP List: ", self.ip_list)
        except:
            print("Input IP.")

        tkWindow = Tk()     # show IP List with Tk
        tkWindow.geometry("220x200")
        tkWindow.title("IP LIST")
        labeltxt = "Number of Collectors: "+str(len(self.ip_list))
        tkLabel = Label(tkWindow, text=labeltxt)
        tkLabel.pack()
        ipListBox = Listbox(tkWindow)
        for i in range(0, len(self.ip_list)):
            ipListBox.insert(i+1, self.ip_list[i])
        ipListBox.pack()
        tkWindow.mainloop()

    def pBtn200_copy_function(self):
        print("200 Copy Pressed")
        self.cmd_list = []  # init. Command List
        self.taskQueue = Queue()  # init. Task Queue
        for i in self.ip_list:
            self.cmd_list.append(make_copy_cmd_list())

        for item in self.cmd_list:
            self.taskQueue.put(item)

        for ip_addr in self.ip_list:
            t = CommandSerize9999(self.taskQueue, ip_addr)
            t.start()
        self.taskQueue.join()
        print("Queue: ", self.taskQueue)

        for i in range(len(self.ip_list)):
            self.taskQueue.put(None)

        for t in self.thread_list:
            t.join()

        print("200 Copy Finish")

    # GCP Tab Function
    def pBtnSetGCPip_function(self):
        print("GCP IP List Set Buttun Pressed.")

        # Load file path with easygui & Load IP List
        filePath = fopen.OpenWinFileExplorer()
        self.gcp_ipListPath.setText(filePath)
        print(self.gcp_ipListPath.text())
        try:
            ip_text = open(self.gcp_ipListPath.text(), "r", encoding='utf8')
        except IOError:
            print("No File.")
            return 0
        self.ip_list = ip_text.read().strip().split("\n")
        print("IP List: ", self.ip_list)
        
        tkWindow = Tk()     # show IP List with Tk
        tkWindow.geometry("220x200")
        tkWindow.title("IP LIST")
        labeltxt = "Number of Collectors: "+str(len(self.ip_list))
        tkLabel = Label(tkWindow, text=labeltxt)
        tkLabel.pack()
        ipListBox = Listbox(tkWindow)
        for i in range(0, len(self.ip_list)):
            ipListBox.insert(i+1, self.ip_list[i])
        ipListBox.pack(side="left", fill="both", expand=True)
        tkWindow.mainloop()
        print("GCP IP List Set Buttun Finish.")
        ip_text.close()

    def pBtnGCP_OpenTargetList_function(self):
        print("GCP Target List Open Pressed")
        filePath = fopen.OpenWinFileExplorer()
        self.gcp_TargetListPath.setText(filePath)
        print(self.gcp_TargetListPath.text())
        try:
            url_list = []
            with open(self.gcp_TargetListPath.text(), "r", encoding='utf8') as f:
                for line in f:
                    url_list.append(line)
            print(url_list)
            tkURLWindow = Tk()     # show IP List with Tk
            tkURLWindow.geometry("320x300")
            tkURLWindow.title("URL LIST")
            labeltxt = "Number of Collectors: "+str(len(url_list))
            tkLabel = Label(tkURLWindow, text=labeltxt)
            tkLabel.pack()
            urlListBox = Listbox(tkURLWindow)
            for i in range(0, len(url_list)):
                urlListBox.insert(i+1, url_list[i])
            urlListBox.pack(side="left", fill="both", expand=True)
            tkURLWindow.mainloop()
        except IOError:
            print("No File.")
            return 0
        print("GCP Target List Open Finish")

    def pBtnGCP_cmd_function(self):
        print("GCP CMD Sending Pressed")
        send_command = self.cmdGCP.text()
        print("Your Command: ",send_command)
        self.cmd_list = []  # init. Command List
        self.taskQueue = Queue()  # init. Task Queue
        for i in self.ip_list:
            self.cmd_list.append(cmd_cmd_list(send_command))

        for item in self.cmd_list:
            self.taskQueue.put(item)

        for ip_addr in self.ip_list:
            # t = CommandSerize9998(self.taskQueue, ip_addr)
            t = CommandSerize9999(self.taskQueue, ip_addr)  # temp modify
            t.start()
        self.taskQueue.join()
        print("Queue: ", self.taskQueue)

        for i in range(len(self.ip_list)):
            self.taskQueue.put(None)

        for t in self.thread_list:
            t.join()
        print("GCP CMD Sending Finish")

    def pBtnGCP_collect_function(self):
        print("GCP Collect Pressed")
        try:
            file_address = self.gcp_TargetListPath.text()
            file_address = file_address.replace("\\", "/")
            file_address = file_address.replace("//", "/")
            file_address = file_address.replace(
                "D:/Tor_CIFS_300/", "/home/jjangga94temp/Tor_CIFS/")
        except:
            print("Input Correct Target List .txt Path")

        if len(self.ip_list) == 0:
            print("IP List is empty.")
        else:
            for i in self.ip_list:
                    self.cmd_list.append(make_tbb_cmd_list(file_address))

            for item in self.cmd_list:
                self.taskQueue.put(item)

            for ip_addr in self.ip_list:
                t = CommandSerize9999(self.taskQueue, ip_addr)
                t.start()
            self.taskQueue.join()
            print("Queue: ", self.taskQueue)

            for i in range(len(self.ip_list)):
                self.taskQueue.put(None)

            for t in self.thread_list:
                t.join()

        print("GCP Collect Finish")

    def pBtnGCP_watch_function(self):
        print("GCP Watch Pressed")
        try:
            file_address = self.gcp_TargetListPath.text()
            file_address = file_address.replace("\\", "/")
            file_address = file_address.replace("//", "/")
            file_address = file_address.replace(
                "D:/Tor_CIFS_300/", "/home/jjangga94temp/Tor_CIFS/")
        except:
            print("Input Correct Target List .txt Path")

        if len(self.ip_list) == 0:
            print("IP List is empty.")
        else:
            for i in self.ip_list:
                self.cmd_list.append(watch_cmd_list(file_address))

            for item in self.cmd_list:
                self.taskQueue.put(item)

            for ip_addr in self.ip_list:
                t = CommandSerize9999(self.taskQueue, ip_addr)
                t.start()
            self.taskQueue.join()
            print("Queue: ", self.taskQueue)

            for i in range(len(self.ip_list)):
                self.taskQueue.put(None)

            for t in self.thread_list:
                t.join()

        print("GCP Watch Finish")

    def pBtnGCP_copy_function(self):
        print("GCP Copy Pressed")
        self.cmd_list = []  # init. Command List
        self.taskQueue = Queue()    # init. Task Queue
        for i in self.ip_list:
            self.cmd_list.append(make_copy_cmd_list())

        for item in self.cmd_list:
            self.taskQueue.put(item)

        for ip_addr in self.ip_list:
            t = CommandSerize9999(self.taskQueue, ip_addr)
            t.start()
        self.taskQueue.join()
        print("Queue: ", self.taskQueue)

        for i in range(len(self.ip_list)):
            self.taskQueue.put(None)

        for t in self.thread_list:
            t.join()

        print("GCP Copy Finish")


# QApplication : Run App
app = QtWidgets.QApplication(sys.argv)

# Create an instance of WindowClass
myWindow = WindowClass()

# Showing the program screen
myWindow.show()

# Code for entering a program into an event loop (which activates the program)
app.exec_()
