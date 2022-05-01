#!/usr/bin/env python3
import camelot
import pandas as pd
import numpy as np
import sys

# usage

# consumer page
# ./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 124-135 consumer-hut.csv

# keyboard page
# ./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 89-95 keyboard-hut.csv

# ~/Downloads/hut1_3_0.pdf
hut = sys.argv[1]
# 123-134
page_range = sys.argv[2]
# output file
out_file = sys.argv[3]
tables = camelot.read_pdf(hut, pages=page_range)

frames = []
for table in tables:
    frames.append(table.df)

ct = pd.concat(frames)

# remove first row (contains column names)
ct = ct.iloc[1:]

# remove irrelevant columns
ct = ct[[0,1]]

# rename columns
ct.rename(columns = {1: 'UsageName'}, inplace=True)
ct.rename(columns = {0: 'UsageID'}, inplace=True)

# remove rows with 'Reserved' string
ct = ct[ct['UsageName'].str.contains('Reserved')==False]

# remove rows 'Usage Name' string
ct = ct[ct['UsageName'].str.contains('Usage Name')==False]

# remove reference numbers e.g. [5]
ct = ct.replace('\[\d*\]', '', regex=True)

# replaced keyboard/keypad words
ct = ct.replace('Keypad', 'KP', regex=True)
ct = ct.replace('Keyboard', 'KB', regex=True)

# remove the and and anything after it
ct = ct.replace('and.*', '', regex=True)

# remove leading/trailing whitespace
ct['UsageName'] = ct['UsageName'].str.strip()

# replace spaces with underscores
ct = ct.replace(' ', '_', regex=True)

# replace slashes with underscores
ct = ct.replace('/', '_', regex=True)

# remove quotes
#ct = ct.replace('"', '', regex=True)

# remove newlines
ct = ct.replace('\n', '', regex=True)

# remove brackets
ct = ct.replace('\(', '', regex=True)
ct = ct.replace('\)', '', regex=True)

# uppercase column
ct['UsageName'] = ct['UsageName'].str.upper()

# convert Usage ID hex to decimal
ct['UsageID'] = ct['UsageID'].apply(int, base=16)

print(ct)
ct.to_csv(out_file, header=False, index=False)
