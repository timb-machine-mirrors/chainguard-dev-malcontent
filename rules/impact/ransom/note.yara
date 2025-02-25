rule ransom_detection: high {
  meta:
    description = "ransomware note"

  strings:
    $not_your_database     = "in your database"
    $not_node              = "NODE_DEBUG_NATIVE"
    $not_private           = "/System/Library/PrivateFrameworks/"
    $not_signed            = "PKCS7_SIGNED"
    $s_already_encrypted   = "already encrypted"
    $s_audit               = "audit of"
    $s_be_decrypted        = "be decrypted"
    $s_blake2b             = "blake2b"
    $s_company             = "company" fullword
    $s_corporate_data      = "corporate data"
    $s_data_recovery       = "data recovery"
    $s_decrypt             = "DECRYPTDIR"
    $s_decrypting          = "Decrypting" fullword
    $s_decryptor           = "decryptor"
    $s_enc_file            = "enc_file"
    $s_enc_stage           = "EncryptionStage"
    $s_encrypt_all         = "encrypt_all"
    $s_encrypt_file        = "encrypt file" nocase
    $s_encrypt_under_file  = "encrypt_file"
    $s_encrypted           = "encrypted" fullword
    $s_encrypted_caps      = "ENCRYPTED"
    $s_encrypting          = "Encrypting" nocase
    $s_entire_network      = "entire network"
    $s_esxcli              = "esxcli"
    $s_esxi                = "esxi"
    $s_EVIDENCE            = "EVIDENCE"
    $s_gained_full         = "gained full access"
    $s_get_in_touch        = "get in touch"
    $s_hi_friends          = "Hi friends"
    $s_iles_txt            = "iles.txt"
    $s_immediate_sale      = "immediate sale"
    $s_incident            = "incident"
    $s_install_tor         = "install TOR"
    $s_insurance           = "insurance"
    $s_is_encrypted        = "s encrypted by"
    $s_key_iv              = "KEY = %s IV = %s"
    $s_lck                 = "%s.lck"
    $s_LEAKAGE             = "LEAKAGE"
    $s_leaks               = "leaks" nocase fullword
    $s_live_chat2          = "live-chat"
    $s_live_chat           = "live chat"
    $s_locker              = "locker" nocase fullword
    $s_lose_access         = "lose access"
    $s_negotiable          = "negotiable"
    $s_negotiation_process = "negotiation process"
    $s_negotiations_open   = "negotiations open"
    $s_negotiators         = "negotiators"
    $s_our_chat            = "our chat"
    $s_our_decryption      = "our decryption"
    $s_our_security        = "our security"
    $s_permanently         = "permanently destroyed"
    $s_ransom              = "ransom" fullword
    $s_recoverfiles        = "recoverfiles"
    $s_recover             = "recover "
    $s_refuse_to_pay       = "refuse to pay"
    $s_remain_silent       = "remain silent"
    $s_restore_docs        = "restore documents" nocase
    $s_restore_my          = "restore-my"
    $s_rypted_f            = "rypted_f"
    $s_this_offer          = "this offer"
    $s_to_board            = "board of directors"
    $s_to_decrypt          = "to decrypt"
    $s_to_my_address       = "to my address"
    $s_unfortunately_your  = "unfortunately your" nocase
    $s_urandom             = "/dev/urandom"
    $s_victim              = "victim"
    $s_vulnerabilities     = "vulnerabilities"
    $s_was_encrypted       = "was encrypted"
    $s_we_stole            = "we stole"
    $s_xlsx                = "xlsx"
    $s_you_can             = "You can request"
    $s_you_decide          = "you decide"
    $s_you_pay             = "you pay "
    $s_your_competitors    = "your competitors"
    $s_your_data           = "your data"
    $s_your_docs           = "your documents" nocase
    $s_your_ench           = "your_enc"
    $s_your_files          = "your files" nocase
    $s_your_network        = "Your network"
    $s_your_security       = "your security"
    $tor_browser           = "TOR Browser" nocase
    $tor_download          = "torproject.org/download"
    $tor_onion             = /\w\.onion\W/

  condition:
    filesize < 20971520 and 6 of ($s_*) and any of ($tor*) and none of ($not*)
}
