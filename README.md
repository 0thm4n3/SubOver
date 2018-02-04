# SubOver

Subover is a Hostile Subdomain Takeover tool designed in Python. From start, it has been aimed with speed and efficiency in mind. Till date, SubOver detects 35 services which is much more than any other tool out there. The tool is multithreaded and hence delivers good speed. It can easily detect and report potential subdomain takeovers that exist. The list of potentially hijackable services is very comprehensive and it is what makes this tool so powerful.

## Installing

You need to have Python 2.7 installed on your machine. The following additional requirements are required - 
- dnspython 
- colorama

Simply clone the repo using `git clone https://github.com/Ice3man543/SubOver.git`.

cd into the repo using `cd SubOver` and run `python subover.py -h`.

## Usage

`python subover.py -l subdomains.txt -o output_takeovers.txt`
- `-l subdomains.txt` is the list of target subdomains. These can be discovered using various tool such as sublist3r or others.
- `-o output_takeovers.txt`is the name of the output file. (Optional & Currently not very well formatted)
- `-t 15` is the number of threads for the tool to use. (Optional)
- `-V` is the switch for showing verbose output. (Optional, Default=False)

## Currently Checked Services

- github
- heroku
- unbounce
- tumblr
- shopify
- instapage
- instapage
- desk
- tictail
- campaignmonitor
- cargocollective
- statuspage
- amazonaws
- cloudfront
- bitbucket
- squarespace
- smartling
- acquia
- fastly
- pantheon
- zendesk
- uservoice
- wpengine
- ghost
- freshdesk
- pingdom
- tilda
- wordpress
- teamwork
- helpjuice
- helpscout
- cargo
- feedpress
- freshdesk
- surge

Count : 35
  
## Screenshot
![tool_in_action](https://raw.githubusercontent.com/Ice3man543/SubOver/master/subover.png)

## FAQ
**Q:** What should my wordlist look like?

**A:** Your wordlist should include a list of subdomains you're checking and should look something like:
```
backend.example.com
something.someone.com
apo-setup.fxc.something.com
```

# TODO

- Add more services :-)
- Improve the tool (There are many things that can be done :-) )

### Development

Want to contribute? Great! 

You can add more services or recommend any changes to the existing ones. Any kind of help is appreciated.

License
----

BSD 2-Clause "Simplified" License


## Contact

Meet me on Twitter: [@Ice3man543](https://twitter.com/ice3man543)
