import datetime
import json
import time
import urllib2

DEFAULT_FILE = 'data.csv'

def GetOandaCandle(startTime, count):
    header = {"X-Accept-Datetime-Format": "UNIX"}
    url = "https://api-fxtrade.oanda.com/v1/candles?instrument=EUR_USD&start=%s&count=%s" % (startTime, count)
    req = urllib2.Request(url, headers=header)
    return json.loads(urllib2.urlopen(req).read())

def fmtCandle(candle):
    t = int(candle['time'][:-6]) # handle us timestamp
    date = datetime.datetime.fromtimestamp(t).strftime('%Y.%m.%d,%H:%M')
    # TODO handle all fields from oanda candle
    return '%s,%s,%s,%s,%s,%s,%s\n' % (date,
                            candle['openBid'],
                            candle['openAsk'],
                            candle['lowBid'],
                            candle['highBid'],
                            candle['closeBid'],
                            candle['volume'])


def writeToFile(data, outfile):
    if 'candles' not in data:
        raise Exception('No candles in this data you fuck face')
    candles = data['candles']
    with open(outfile, 'w') as f:
        for candle in candles:
            fmtdCandle = fmtCandle(candle)
            f.write(fmtdCandle)

def main():
    # TODO allow user specify a year or month or whatever this goes and gets a years or month or whatever worth of data and writes it to a csv (that eventually the user will specify)
    data = GetOandaCandle(1471365995, 3)
    writeToFile(data, DEFAULT_FILE)

main()
