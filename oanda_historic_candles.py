import argparse
import datetime
import json
import time
import urllib2

DEFAULT_FILE = 'data.csv'
MONTH_SEC = 60 * 60 * 24 * 30

def WriteMonthOfData(month, year, f, startTime=None):
    if not startTime:
        if not year:
            raise 'if you dont provide startTime you must provide a year'
        if not month:
            month = 1 # default to jan
        month = '0' + str(month)
        date = "%s-%s-01" % (year, month)
        dt = datetime.datetime.strptime(date, '%Y-%m-%d')
        start = int(time.mktime(dt.timetuple()))
    else:
        start = startTime
    endTime = start + MONTH_SEC

    t = start
    while t < endTime:
        # get finer grain number of records
        # 5000 ensures a month of data is retrieved without
        # getting our requests throttled but means we
        # can overshoot by a lot
        data = GetOandaCandle(start, 5000)
        writeToFile(data, f)
        t = int(data['candles'][-1]['time'][:-6])
    return data['candles'][-1]['time'][:-6]

# TODO do something to avoid hitting the rate limit
def WriteYearOfData(year, f):
    lastTime = None
    for i in xrange(1, 13):
        if lastTime:
            lastTime = WriteMonthOfData(i, year, f, lastTime)
        else:
            lastTime = WriteMonthOfData(i, year, f)

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
    with open(outfile, 'a') as f:
        for candle in candles:
            fmtdCandle = fmtCandle(candle)
            f.write(fmtdCandle)

def main():
    parser = argparse.ArgumentParser(description='Specify time period to retrieve data from oanda')
    parser.add_argument('--year', required=True, type=int, help='year to get data for')
    parser.add_argument('--month', required=False, type=int, help='month to get data for')
    parser.add_argument('--outfile', required=False, help='file to write data to')

    args = parser.parse_args()

    if args.month:
        WriteMonthOfData(args.month, args.year, DEFAULT_FILE)
    else:
        WriteYearOfData(args.year, DEFAULT_FILE)

main()
