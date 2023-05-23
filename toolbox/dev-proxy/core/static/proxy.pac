function FindProxyForURL(url, host) {
  url = url.toLowerCase();
  host = host.toLowerCase();

  if (url.startsWith("http:")) {
    return "PROXY 127.0.0.1:1234";
  }

//   if (
//     isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
//     isInNet(dnsResolve(host), "172.16.0.0", "255.240.0.0") ||
//     isInNet(dnsResolve(host), "172.1.1.0", "255.255.255.0") ||
//     isInNet(dnsResolve(host), "192.168.0.0", "255.255.0.0") ||
//     isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0")
//   ) {
//     return "DIRECT";
//   }

  if (
    shExpMatch(url, "*google.com*") ||
    shExpMatch(url, "*google.co*") ||
    shExpMatch(url, "*gmail.com*") ||
    shExpMatch(url, "*google.dev*") ||
    shExpMatch(url, "*twitter.com*") ||
    shExpMatch(url, "*github.com*") ||
    shExpMatch(url, "*youtube.com*") ||
    shExpMatch(url, "*wikipedia.org*")
  ) {
    return "PROXY 127.0.0.1:7890";
  }

  return "DIRECT";
}
