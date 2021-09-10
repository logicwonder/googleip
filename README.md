# googleip
A command line utility to keeps track of Google public IP ranges. 

This tool helps you to track changes made (additions/deletions) in the Google public IP allow list to configure your firewall rules on a periodic basis.

Read the about the requirement here:
https://firebase.google.com/docs/cloud-messaging/concept-options#messaging-ports-and-your-firewall
> If your organization has a firewall to restrict traffic to or from the Internet, you need to configure it to allow mobile devices to connect with FCM in order for devices on your network to receive messages. FCM typically uses port 5228, but it sometimes uses 443, 5229, and 5230. For devices connecting on your network, FCM doesn't provide specific IPs because our IP range changes too frequently and your firewall rules could get out of date, impacting your users' experience. Ideally, allowlist ports 5228-5230 & 443 with no IP restrictions. However, if you must have an IP restriction, you should allowlist all of the IP addresses listed in goog.json. This large list is updated regularly, and you are recommended to update your rules on a monthly basis. Problems caused by firewall IP restrictions are often intermittent and difficult to diagnose.See 

## Usage
1. Build the tool
```
make build 
cd bin/
````

2. Run the following initially to download the latest IP allow list from Google.
```
googleip -refresh
```

3. Run the following to dump the IP4/IPV6 allow list.

```
googleip -dump=ipv4

googleip -dump=ipv6

```

4. Periodically checking for changes in IP allow list.
```
googleip -refresh
```

**Note**:- The current and archived IP allow lists are stored in scribbledb/goog/ folder