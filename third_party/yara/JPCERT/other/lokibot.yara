rule malware_Lokibot_strings {
          meta:
            description = "detect Lokibot in memory"
            author = "JPCERT/CC Incident Response Group"
            rule_usage = "memory scan"
            reference = "internal research"


          strings:
            $des3 = { 68 03 66 00 00 }
            $param = "MAC=%02X%02X%02XINSTALL=%08X%08X"
            $string = { 2d 00 75 00 00 00 46 75 63 6b 61 76 2e 72 75 00 00}

          condition:
            all of them
}
