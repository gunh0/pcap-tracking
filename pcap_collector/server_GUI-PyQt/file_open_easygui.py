# file_open_easygui.py
# fileopenbox: fileopenbox returns the name of a file

import easygui

def OpenWinFileExplorer():
    multiSearchFilePath = easygui.fileopenbox()
    return multiSearchFilePath