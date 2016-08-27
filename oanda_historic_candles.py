import argparse
import datetime
from datetime import date
import json
import time
import urllib2
from calendar import monthrange

DEFAULT_FILE = 'data.csv'
MONTH_SEC = 60 * 60 * 24 * 30
DAY_SEC = 60 * 60 * 24
TOO_LOW = ['S5', 'S10', 'S15']
LOW_GRAN = ['S30', 'M1', 'M2', 'M3', 'M4', 'M5']

# Reference: http://developer.oanda.com/rest-live/rates/
# Example: python oanda_historic_candles.py  --year 2012 --granularity M10

def WriteMonthOfData(month, year, granularity, f):
    if not year:
        raise 'if you dont provide startTime you must provide a year'
    if not month:
        month = 1 # default to jan

    if month > 9:
        month = str(month)
    else:
        month = '0' + str(month)
    date = "%s-%s-01" % (year, month)
    dt = datetime.datetime.strptime(date, '%Y-%m-%d')
    start = int(time.mktime(dt.timetuple()))

    end = start + MONTH_SEC

    data = GetOandaCandle(start, end, granularity)
    writeToFile(data, f)
    return data['candles'][-1]['time'][:-6]

def WriteGranularMonth(month, year, granularity, f):
    if month == 12:
        days = (date(year+1, 1, 1) - date(year, month, 1)).days
    else:
        days = (date(year, month+1, 1) - date(year, month, 1)).days

    for i in xrange(1, days + 1):
        WriteDayOfData(i, month, year, granularity, f)


def WriteDayOfData(day, month, year, granularity, f):
    if not year:
        raise 'if you dont provide startTime you must provide a year'
    if not month:
        month = 1 # default to jan

    if month > 9:
        month = str(month)
    else:
        month = '0' + str(month)
    date = "%s-%s-%s" % (year, month, day)
    dt = datetime.datetime.strptime(date, '%Y-%m-%d')
    start = int(time.mktime(dt.timetuple()))

    end = start + DAY_SEC

    data = GetOandaCandle(start, end, granularity)
    writeToFile(data, f)
    return data['candles'][-1]['time'][:-6]

def WriteYearOfData(year, granularity, f):
    lastTime = None
    for i in xrange(1, 13):
        WriteMonthOfData(i, year, granularity, f)

def WriteGranularYear(year, granularity, f):
    lastTime = None
    for i in xrange(1, 13):
        WriteGranularMonth(i, year, granularity, f)

def GetOandaCandle(startTime, endTime, granularity):
    ret = SafeGetOandaCandle(startTime, endTime, granularity)
    if not ret:
        time.sleep(120)
        return SafeGetOandaCandle(startTime, endTime, granularity)
    return ret

def SafeGetOandaCandle(startTime, endTime, granularity):
    header = {"X-Accept-Datetime-Format": "UNIX"}
    url = "https://api-fxtrade.oanda.com/v1/candles?instrument=EUR_USD&start=%s&end=%s&granularity=%s" % (startTime, endTime, granularity)
    req = urllib2.Request(url, headers=header)
    try:
        resp = urllib2.urlopen(req)
        return json.loads(resp.read())
    except urllib2.HTTPError as e:
        if e.code == 429:
            print 'you are being rate limited. Im going to sleep 2 minutes and try again.'
            return None


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
    parser.add_argument('--granularity', required=True, type=str, help='granularity of data ie M1 = 1 minute')
    parser.add_argument('--outfile', required=False, help='file to write data to')

    args = parser.parse_args()

    if args.granularity in TOO_LOW:
        raise 'We cant go that low in granularity m8'

    if args.month:
        if args.granularity in LOW_GRAN:
            WriteGranularMonth(args.month, args.year, args.granularity, DEFAULT_FILE)
        else:
            WriteMonthOfData(args.month, args.year, args.granularity, DEFAULT_FILE)
    else:
        if args.granularity in LOW_GRAN:
            print 'it will take like 10 minutes to get a year of granular data'
            print 'im still going to get that data but you should know it will take long af'
            WriteGranularYear(args.year, args.granularity, DEFAULT_FILE)
        else:
            WriteYearOfData(args.year, args.granularity, DEFAULT_FILE)

main()
