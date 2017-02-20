#!/usr/bin/env python
import os
import sys
import time
import subprocess
import logging

import pygle.config
pygle.config.user = os.environ['wigle_api_user']
pygle.config.key = os.environ['wigle_api_key']

import pygle.network
import pygle.requests
import pygle.profile


def setup_logging():
    try:
        import http.client as http_client
    except ImportError:
        # Python 2
        import httplib as http_client
    http_client.HTTPConnection.debuglevel = 1

    logging.basicConfig()
    logging.getLogger().setLevel(logging.DEBUG)
    requests_log = logging.getLogger("requests.packages.urllib3")
    requests_log.setLevel(logging.DEBUG)
    requests_log.propagate = True


def get_bssids():
    bssid_output = subprocess.check_output(
        "iw wlan0 scan | grep ^BSS | cut -f 1 -d '(' | awk '{ print $2}'",
        stderr=subprocess.STDOUT,
        shell=True)
    if "command failed: Operation not permitted (-1)" in bssid_output:
        print(bssid_output)
        sys.exit(1)
    return set(bssid_output.strip().split("\n"))


def get_lat_long(bssids):
    bssid_data = {}
    best_bssid = ""
#>>> pygle.network.detail(netid="F8:D1:11:AD:99:BE")
#{u'wifi': True, u'gsm': False, u'results': [{u'comment': None, u'firsttime': u'2013-05-18T06:02:01.000Z', u'ssid': u'?? ? =?????', u'lastupdt': u'2013-05-21T22:44:24.000Z', u'trilong': -122.75737762, u'netid': u'F8:D1:11:AD:99:BE', u'freenet': u'?', u'paynet': u'?', u'userfound': None, u'wep': u'2', u'bcninterval': 0, u'dhcp': u'?', u'trilat': 47.059612270000002, u'qos': 0, u'transid': u'20130521-00351', u'lasttime': u'2013-05-22T00:44:19.000Z', u'type': u'infra', u'locationData': [{u'noise': 0.0, u'ssid': u'?? ? =?????', u'name': None, u'lastupdt': u'2013-05-21T22:44:19.000Z', u'signal': -93.0, u'netId': u'273576828443070', u'longitude': -122.75737762, u'month': u'201305', u'wep': u'2', u'snr': 0.0, u'time': u'2013-05-18T06:02:01.000Z', u'latitude': 47.059612270000002, u'alt': 50, u'encryptionValue': u'WPA2', u'accuracy': 5.0999999999999996}], u'channel': 6, u'name': None}], u'success': True, u'cdma': False}
    for bssid in bssids:
        response = pygle.network.detail(netid=bssid)
        print(response)
        if not response.get('success', False):
            print("Query for {} failed!".format(bssid))
            sys.exit(1)
        results = response['results']
        print(results)
        if len(results) != 1:
            print("Warning, multiple results for {}".format(bssid))
        bssid_data[bssid] = results[0]
        lastupdt = bssid_data[bssid]['lastupdt']
        if lastupdt > bssid_data[best_bssid]['lastupdt']:
            print("{} is the best so far updated {}".format(bssid, lastupdt))
            best_bssid = bssid
        time.sleep(60)
    print("Best (most recent) result is:")
    print(bssid_data[best_bssid])


if __name__ == "__main__":
    bssids = get_bssids()
    print("Got bssids: {}".format(",".join(bssids)))
    lat, lon = get_lat_long(bssids)
