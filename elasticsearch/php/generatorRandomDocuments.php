<?php

class generatorRandomDocuments {
  
  protected $licenses;
  protected $licenses_numbers;
  
  protected $main_copyright_owner;
  protected $main_copyright_owner_numbers;
  
  protected $repo_owner;
  protected $repo_owner_numbers;
  
  protected $maintainance_type;
  protected $maintainance_type_numbers;
  
  protected $technical_contacts;

  protected $platforms;
  protected $platforms_numbers;

  protected $descriptions;

  protected $scope;
  protected $scope_numbers;

  protected $patype;
  protected $patype_numbers;

  protected $usedby;
  protected $usedby_numbers;

  protected $tags;
  protected $tags_numbers;

  protected $category;
  protected $category_numbers;

  protected $dependencies;
  protected $dependencies_numbers;

  protected $dependencies_hardware;
  protected $dependencies_hardware_numbers;

  protected $metadata_repo;
  protected $metadata_repo_numbers;

  protected $development_status;
  protected $development_status_numbers;

  protected $software_type;
  protected $software_type_numbers;

  protected $ISO_639_3;
  protected $ISO_639_3_numbers;

  protected $mime_types;
  protected $mime_types_numbers;

  protected $ecosistemi;
  protected $ecosistemi_numbers;

  public function __construct() {
    $this->licenses = [
      "0BSD",                                 // BSD Zero Clause License
      "AAL",                                  // Attribution Assurance License
      "Abstyles",                             // Abstyles License
      "Adobe-2006",                           // Adobe Systems Incorporated Source Code License Agreement
      "Adobe-Glyph",                          // Adobe Glyph List License
      "ADSL",                                 // Amazon Digital Services License
      "AFL-1.1",                              // Academic Free License v1.1
      "AFL-1.2",                              // Academic Free License v1.2
      "AFL-2.0",                              // Academic Free License v2.0
      "AFL-2.1",                              // Academic Free License v2.1
      "AFL-3.0",                              // Academic Free License v3.0
      "Afmparse",                             // Afmparse License
      "AGPL-1.0-only",                        // Affero General Public License v1.0 only
      "AGPL-1.0-or-later",                    // Affero General Public License v1.0 or later
      "AGPL-3.0-only",                        // GNU Affero General Public License v3.0 only
      "AGPL-3.0-or-later",                    // GNU Affero General Public License v3.0 or later
      "Aladdin",                              // Aladdin Free Public License
      "AMDPLPA",                              // AMD's plpa_map.c License
      "AML",                                  // Apple MIT License
      "AMPAS",                                // Academy of Motion Picture Arts and Sciences BSD
      "ANTLR-PD",                             // ANTLR Software Rights Notice
      "Apache-1.0",                           // Apache License 1.0
      "Apache-1.1",                           // Apache License 1.1
      "Apache-2.0",                           // Apache License 2.0
      "APAFML",                               // Adobe Postscript AFM License
      "APL-1.0",                              // Adaptive Public License 1.0
      "APSL-1.0",                             // Apple Public Source License 1.0
      "APSL-1.1",                             // Apple Public Source License 1.1
      "APSL-1.2",                             // Apple Public Source License 1.2
      "APSL-2.0",                             // Apple Public Source License 2.0
      "Artistic-1.0",                         // Artistic License 1.0
      "Artistic-1.0-cl8",                     // Artistic License 1.0 w/clause 8
      "Artistic-1.0-Perl",                    // Artistic License 1.0 (Perl)
      "Artistic-2.0",                         // Artistic License 2.0
      "Bahyph",                               // Bahyph License
      "Barr",                                 // Barr License
      "Beerware",                             // Beerware License
      "BitTorrent-1.0",                       // BitTorrent Open Source License v1.0
      "BitTorrent-1.1",                       // BitTorrent Open Source License v1.1
      "Borceux",                              // Borceux license
      "BSD-1-Clause",                         // BSD 1-Clause License
      "BSD-2-Clause",                         // BSD 2-Clause "Simplified" License
      "BSD-2-Clause-FreeBSD",                 // BSD 2-Clause FreeBSD License
      "BSD-2-Clause-NetBSD",                  // BSD 2-Clause NetBSD License
      "BSD-2-Clause-Patent",                  // BSD-2-Clause Plus Patent License
      "BSD-3-Clause",                         // BSD 3-Clause "New" or "Revised" License
      "BSD-3-Clause-Attribution",             // BSD with attribution
      "BSD-3-Clause-Clear",                   // BSD 3-Clause Clear License
      "BSD-3-Clause-LBNL",                    // Lawrence Berkeley National Labs BSD variant license
      "BSD-3-Clause-No-Nuclear-License",      // BSD 3-Clause No Nuclear License
      "BSD-3-Clause-No-Nuclear-License-2014", // BSD 3-Clause No Nuclear License 2014
      "BSD-3-Clause-No-Nuclear-Warranty",     // BSD 3-Clause No Nuclear Warranty
      "BSD-4-Clause",                         // BSD 4-Clause "Original" or "Old" License
      "BSD-4-Clause-UC",                      // BSD-4-Clause (University of California-Specific)
      "BSD-Protection",                       // BSD Protection License
      "BSD-Source-Code",                      // BSD Source Code Attribution
      "BSL-1.0",                              // Boost Software License 1.0
      "bzip2-1.0.5",                          // bzip2 and libbzip2 License v1.0.5
      "bzip2-1.0.6",                          // bzip2 and libbzip2 License v1.0.6
      "Caldera",                              // Caldera License
      "CATOSL-1.1",                           // Computer Associates Trusted Open Source License 1.1
      "CC-BY-1.0",                            // Creative Commons Attribution 1.0 Generic
      "CC-BY-2.0",                            // Creative Commons Attribution 2.0 Generic
      "CC-BY-2.5",                            // Creative Commons Attribution 2.5 Generic
      "CC-BY-3.0",                            // Creative Commons Attribution 3.0 Unported
      "CC-BY-4.0",                            // Creative Commons Attribution 4.0 International
      "CC-BY-NC-1.0",                         // Creative Commons Attribution Non Commercial 1.0 Generic
      "CC-BY-NC-2.0",                         // Creative Commons Attribution Non Commercial 2.0 Generic
      "CC-BY-NC-2.5",                         // Creative Commons Attribution Non Commercial 2.5 Generic
      "CC-BY-NC-3.0",                         // Creative Commons Attribution Non Commercial 3.0 Unported
      "CC-BY-NC-4.0",                         // Creative Commons Attribution Non Commercial 4.0 International
      "CC-BY-NC-ND-1.0",                      // Creative Commons Attribution Non Commercial No Derivatives 1.0 Generic
      "CC-BY-NC-ND-2.0",                      // Creative Commons Attribution Non Commercial No Derivatives 2.0 Generic
      "CC-BY-NC-ND-2.5",                      // Creative Commons Attribution Non Commercial No Derivatives 2.5 Generic
      "CC-BY-NC-ND-3.0",                      // Creative Commons Attribution Non Commercial No Derivatives 3.0 Unported
      "CC-BY-NC-ND-4.0",                      // Creative Commons Attribution Non Commercial No Derivatives 4.0 International
      "CC-BY-NC-SA-1.0",                      // Creative Commons Attribution Non Commercial Share Alike 1.0 Generic
      "CC-BY-NC-SA-2.0",                      // Creative Commons Attribution Non Commercial Share Alike 2.0 Generic
      "CC-BY-NC-SA-2.5",                      // Creative Commons Attribution Non Commercial Share Alike 2.5 Generic
      "CC-BY-NC-SA-3.0",                      // Creative Commons Attribution Non Commercial Share Alike 3.0 Unported
      "CC-BY-NC-SA-4.0",                      // Creative Commons Attribution Non Commercial Share Alike 4.0 International
      "CC-BY-ND-1.0",                         // Creative Commons Attribution No Derivatives 1.0 Generic
      "CC-BY-ND-2.0",                         // Creative Commons Attribution No Derivatives 2.0 Generic
      "CC-BY-ND-2.5",                         // Creative Commons Attribution No Derivatives 2.5 Generic
      "CC-BY-ND-3.0",                         // Creative Commons Attribution No Derivatives 3.0 Unported
      "CC-BY-ND-4.0",                         // Creative Commons Attribution No Derivatives 4.0 International
      "CC-BY-SA-1.0",                         // Creative Commons Attribution Share Alike 1.0 Generic
      "CC-BY-SA-2.0",                         // Creative Commons Attribution Share Alike 2.0 Generic
      "CC-BY-SA-2.5",                         // Creative Commons Attribution Share Alike 2.5 Generic
      "CC-BY-SA-3.0",                         // Creative Commons Attribution Share Alike 3.0 Unported
      "CC-BY-SA-4.0",                         // Creative Commons Attribution Share Alike 4.0 International
      "CC0-1.0",                              // Creative Commons Zero v1.0 Universal
      "CDDL-1.0",                             // Common Development and Distribution License 1.0
      "CDDL-1.1",                             // Common Development and Distribution License 1.1
      "CDLA-Permissive-1.0",                  // Community Data License Agreement Permissive 1.0
      "CDLA-Sharing-1.0",                     // Community Data License Agreement Sharing 1.0
      "CECILL-1.0",                           // CeCILL Free Software License Agreement v1.0
      "CECILL-1.1",                           // CeCILL Free Software License Agreement v1.1
      "CECILL-2.0",                           // CeCILL Free Software License Agreement v2.0
      "CECILL-2.1",                           // CeCILL Free Software License Agreement v2.1
      "CECILL-B",                             // CeCILL-B Free Software License Agreement
      "CECILL-C",                             // CeCILL-C Free Software License Agreement
      "ClArtistic",                           // Clarified Artistic License
      "CNRI-Jython",                          // CNRI Jython License
      "CNRI-Python",                          // CNRI Python License
      "CNRI-Python-GPL-Compatible",           // CNRI Python Open Source GPL Compatible License Agreement
      "Condor-1.1",                           // Condor Public License v1.1
      "CPAL-1.0",                             // Common Public Attribution License 1.0
      "CPL-1.0",                              // Common Public License 1.0
      "CPOL-1.02",                            // Code Project Open License 1.02
      "Crossword",                            // Crossword License
      "CrystalStacker",                       // CrystalStacker License
      "CUA-OPL-1.0",                          // CUA Office Public License v1.0
      "Cube",                                 // Cube License
      "curl",                                 // curl License
      "D-FSL-1.0",                            // Deutsche Freie Software Lizenz
      "diffmark",                             // diffmark license
      "DOC",                                  // DOC License
      "Dotseqn",                              // Dotseqn License
      "DSDP",                                 // DSDP License
      "dvipdfm",                              // dvipdfm License
      "ECL-1.0",                              // Educational Community License v1.0
      "ECL-2.0",                              // Educational Community License v2.0
      "EFL-1.0",                              // Eiffel Forum License v1.0
      "EFL-2.0",                              // Eiffel Forum License v2.0
      "eGenix",                               // eGenix.com Public License 1.1.0
      "Entessa",                              // Entessa Public License v1.0
      "EPL-1.0",                              // Eclipse Public License 1.0
      "EPL-2.0",                              // Eclipse Public License 2.0
      "ErlPL-1.1",                            // Erlang Public License v1.1
      "EUDatagrid",                           // EU DataGrid Software License
      "EUPL-1.0",                             // European Union Public License 1.0
      "EUPL-1.1",                             // European Union Public License 1.1
      "EUPL-1.2",                             // European Union Public License 1.2
      "Eurosym",                              // Eurosym License
      "Fair",                                 // Fair License
      "Frameworx-1.0",                        // Frameworx Open License 1.0
      "FreeImage",                            // FreeImage Public License v1.0
      "FSFAP",                                // FSF All Permissive License
      "FSFUL",                                // FSF Unlimited License
      "FSFULLR",                              // FSF Unlimited License (with License Retention)
      "FTL",                                  // Freetype Project License
      "GFDL-1.1-only",                        // GNU Free Documentation License v1.1 only
      "GFDL-1.1-or-later",                    // GNU Free Documentation License v1.1 or later
      "GFDL-1.2-only",                        // GNU Free Documentation License v1.2 only
      "GFDL-1.2-or-later",                    // GNU Free Documentation License v1.2 or later
      "GFDL-1.3-only",                        // GNU Free Documentation License v1.3 only
      "GFDL-1.3-or-later",                    // GNU Free Documentation License v1.3 or later
      "Giftware",                             // Giftware License
      "GL2PS",                                // GL2PS License
      "Glide",                                // 3dfx Glide License
      "Glulxe",                               // Glulxe License
      "gnuplot",                              // gnuplot License
      "GPL-1.0-only",                         // GNU General Public License v1.0 only
      "GPL-1.0-or-later",                     // GNU General Public License v1.0 or later
      "GPL-2.0-only",                         // GNU General Public License v2.0 only
      "GPL-2.0-or-later",                     // GNU General Public License v2.0 or later
      "GPL-3.0-only",                         // GNU General Public License v3.0 only
      "GPL-3.0-or-later",                     // GNU General Public License v3.0 or later
      "gSOAP-1.3b",                           // gSOAP Public License v1.3b
      "HaskellReport",                        // Haskell Language Report License
      "HPND",                                 // Historical Permission Notice and Disclaimer
      "IBM-pibs",                             // IBM PowerPC Initialization and Boot Software
      "ICU",                                  // ICU License
      "IJG",                                  // Independent JPEG Group License
      "ImageMagick",                          // ImageMagick License
      "iMatix",                               // iMatix Standard Function Library Agreement
      "Imlib2",                               // Imlib2 License
      "Info-ZIP",                             // Info-ZIP License
      "Intel",                                // Intel Open Source License
      "Intel-ACPI",                           // Intel ACPI Software License Agreement
      "Interbase-1.0",                        // Interbase Public License v1.0
      "IPA",                                  // IPA Font License
      "IPL-1.0",                              // IBM Public License v1.0
      "ISC",                                  // ISC License
      "JasPer-2.0",                           // JasPer License
      "JSON",                                 // JSON License
      "LAL-1.2",                              // Licence Art Libre 1.2
      "LAL-1.3",                              // Licence Art Libre 1.3
      "Latex2e",                              // Latex2e License
      "Leptonica",                            // Leptonica License
      "LGPL-2.0-only",                        // GNU Library General Public License v2 only
      "LGPL-2.0-or-later",                    // GNU Library General Public License v2 or later
      "LGPL-2.1-only",                        // GNU Lesser General Public License v2.1 only
      "LGPL-2.1-or-later",                    // GNU Lesser General Public License v2.1 or later
      "LGPL-3.0-only",                        // GNU Lesser General Public License v3.0 only
      "LGPL-3.0-or-later",                    // GNU Lesser General Public License v3.0 or later
      "LGPLLR",                               // Lesser General Public License For Linguistic Resources
      "Libpng",                               // libpng License
      "libtiff",                              // libtiff License
      "LiLiQ-P-1.1",                          // Licence Libre du Québec – Permissive version 1.1
      "LiLiQ-R-1.1",                          // Licence Libre du Québec – Réciprocité version 1.1
      "LiLiQ-Rplus-1.1",                      // Licence Libre du Québec – Réciprocité forte version 1.1
      "Linux-OpenIB",                         // Linux Kernel Variant of OpenIB.org license
      "LPL-1.0",                              // Lucent Public License Version 1.0
      "LPL-1.02",                             // Lucent Public License v1.02
      "LPPL-1.0",                             // LaTeX Project Public License v1.0
      "LPPL-1.1",                             // LaTeX Project Public License v1.1
      "LPPL-1.2",                             // LaTeX Project Public License v1.2
      "LPPL-1.3a",                            // LaTeX Project Public License v1.3a
      "LPPL-1.3c",                            // LaTeX Project Public License v1.3c
      "MakeIndex",                            // MakeIndex License
      "MirOS",                                // MirOS License
      "MIT",                                  // MIT License
      "MIT-0",                                // MIT No Attribution
      "MIT-advertising",                      // Enlightenment License (e16)
      "MIT-CMU",                              // CMU License
      "MIT-enna",                             // enna License
      "MIT-feh",                              // feh License
      "MITNFA",                               // MIT +no-false-attribs license
      "Motosoto",                             // Motosoto License
      "mpich2",                               // mpich2 License
      "MPL-1.0",                              // Mozilla Public License 1.0
      "MPL-1.1",                              // Mozilla Public License 1.1
      "MPL-2.0",                              // Mozilla Public License 2.0
      "MPL-2.0-no-copyleft-exception",        // Mozilla Public License 2.0 (no copyleft exception)
      "MS-PL",                                // Microsoft Public License
      "MS-RL",                                // Microsoft Reciprocal License
      "MTLL",                                 // Matrix Template Library License
      "Multics",                              // Multics License
      "Mup",                                  // Mup License
      "NASA-1.3",                             // NASA Open Source Agreement 1.3
      "Naumen",                               // Naumen Public License
      "NBPL-1.0",                             // Net Boolean Public License v1
      "NCSA",                                 // University of Illinois/NCSA Open Source License
      "Net-SNMP",                             // Net-SNMP License
      "NetCDF",                               // NetCDF license
      "Newsletr",                             // Newsletr License
      "NGPL",                                 // Nethack General Public License
      "NLOD-1.0",                             // Norwegian Licence for Open Government Data
      "NLPL",                                 // No Limit Public License
      "Nokia",                                // Nokia Open Source License
      "NOSL",                                 // Netizen Open Source License
      "Noweb",                                // Noweb License
      "NPL-1.0",                              // Netscape Public License v1.0
      "NPL-1.1",                              // Netscape Public License v1.1
      "NPOSL-3.0",                            // Non-Profit Open Software License 3.0
      "NRL",                                  // NRL License
      "NTP",                                  // NTP License
      "OCCT-PL",                              // Open CASCADE Technology Public License
      "OCLC-2.0",                             // OCLC Research Public License 2.0
      "ODbL-1.0",                             // ODC Open Database License v1.0
      "OFL-1.0",                              // SIL Open Font License 1.0
      "OFL-1.1",                              // SIL Open Font License 1.1
      "OGTSL",                                // Open Group Test Suite License
      "OLDAP-1.1",                            // Open LDAP Public License v1.1
      "OLDAP-1.2",                            // Open LDAP Public License v1.2
      "OLDAP-1.3",                            // Open LDAP Public License v1.3
      "OLDAP-1.4",                            // Open LDAP Public License v1.4
      "OLDAP-2.0",                            // Open LDAP Public License v2.0 (or possibly 2.0A and 2.0B)
      "OLDAP-2.0.1",                          // Open LDAP Public License v2.0.1
      "OLDAP-2.1",                            // Open LDAP Public License v2.1
      "OLDAP-2.2",                            // Open LDAP Public License v2.2
      "OLDAP-2.2.1",                          // Open LDAP Public License v2.2.1
      "OLDAP-2.2.2",                          // Open LDAP Public License 2.2.2
      "OLDAP-2.3",                            // Open LDAP Public License v2.3
      "OLDAP-2.4",                            // Open LDAP Public License v2.4
      "OLDAP-2.5",                            // Open LDAP Public License v2.5
      "OLDAP-2.6",                            // Open LDAP Public License v2.6
      "OLDAP-2.7",                            // Open LDAP Public License v2.7
      "OLDAP-2.8",                            // Open LDAP Public License v2.8
      "OML",                                  // Open Market License
      "OpenSSL",                              // OpenSSL License
      "OPL-1.0",                              // Open Public License v1.0
      "OSET-PL-2.1",                          // OSET Public License version 2.1
      "OSL-1.0",                              // Open Software License 1.0
      "OSL-1.1",                              // Open Software License 1.1
      "OSL-2.0",                              // Open Software License 2.0
      "OSL-2.1",                              // Open Software License 2.1
      "OSL-3.0",                              // Open Software License 3.0
      "PDDL-1.0",                             // ODC Public Domain Dedication & License 1.0
      "PHP-3.0",                              // PHP License v3.0
      "PHP-3.01",                             // PHP License v3.01
      "Plexus",                               // Plexus Classworlds License
      "PostgreSQL",                           // PostgreSQL License
      "psfrag",                               // psfrag License
      "psutils",                              // psutils License
      "Python-2.0",                           // Python License 2.0
      "Qhull",                                // Qhull License
      "QPL-1.0",                              // Q Public License 1.0
      "Rdisc",                                // Rdisc License
      "RHeCos-1.1",                           // Red Hat eCos Public License v1.1
      "RPL-1.1",                              // Reciprocal Public License 1.1
      "RPL-1.5",                              // Reciprocal Public License 1.5
      "RPSL-1.0",                             // RealNetworks Public Source License v1.0
      "RSA-MD",                               // RSA Message-Digest License
      "RSCPL",                                // Ricoh Source Code Public License
      "Ruby",                                 // Ruby License
      "SAX-PD",                               // Sax Public Domain Notice
      "Saxpath",                              // Saxpath License
      "SCEA",                                 // SCEA Shared Source License
      "Sendmail",                             // Sendmail License
      "SGI-B-1.0",                            // SGI Free Software License B v1.0
      "SGI-B-1.1",                            // SGI Free Software License B v1.1
      "SGI-B-2.0",                            // SGI Free Software License B v2.0
      "SimPL-2.0",                            // Simple Public License 2.0
      "SISSL",                                // Sun Industry Standards Source License v1.1
      "SISSL-1.2",                            // Sun Industry Standards Source License v1.2
      "Sleepycat",                            // Sleepycat License
      "SMLNJ",                                // Standard ML of New Jersey License
      "SMPPL",                                // Secure Messaging Protocol Public License
      "SNIA",                                 // SNIA Public License 1.1
      "Spencer-86",                           // Spencer License 86
      "Spencer-94",                           // Spencer License 94
      "Spencer-99",                           // Spencer License 99
      "SPL-1.0",                              // Sun Public License v1.0
      "SugarCRM-1.1.3",                       // SugarCRM Public License v1.1.3
      "SWL",                                  // Scheme Widget Library (SWL) Software License Agreement
      "TCL",                                  // TCL/TK License
      "TCP-wrappers",                         // TCP Wrappers License
      "TMate",                                // TMate Open Source License
      "TORQUE-1.1",                           // TORQUE v2.5+ Software License v1.1
      "TOSL",                                 // Trusster Open Source License
      "Unicode-DFS-2015",                     // Unicode License Agreement - Data Files and Software (2015)
      "Unicode-DFS-2016",                     // Unicode License Agreement - Data Files and Software (2016)
      "Unicode-TOU",                          // Unicode Terms of Use
      "Unlicense",                            // The Unlicense
      "UPL-1.0",                              // Universal Permissive License v1.0
      "Vim",                                  // Vim License
      "VOSTROM",                              // VOSTROM Public License for Open Source
      "VSL-1.0",                              // Vovida Software License v1.0
      "W3C",                                  // W3C Software Notice and License (2002-12-31)
      "W3C-19980720",                         // W3C Software Notice and License (1998-07-20)
      "W3C-20150513",                         // W3C Software Notice and Document License (2015-05-13)
      "Watcom-1.0",                           // Sybase Open Watcom Public License 1.0
      "Wsuipa",                               // Wsuipa License
      "WTFPL",                                // Do What The F*ck You Want To Public License
      "X11",                                  // X11 License
      "Xerox",                                // Xerox License
      "XFree86-1.1",                          // XFree86 License 1.1
      "xinetd",                               // xinetd License
      "Xnet",                                 // X.Net License
      "xpp",                                  // XPP License
      "XSkat",                                // XSkat License
      "YPL-1.0",                              // Yahoo! Public License v1.0
      "YPL-1.1",                              // Yahoo! Public License v1.1
      "Zed",                                  // Zed License
      "Zend-2.0",                             // Zend License v2.0
      "Zimbra-1.3",                           // Zimbra Public License v1.3
      "Zimbra-1.4",                           // Zimbra Public License v1.4
      "Zlib",                                 // zlib License
      "zlib-acknowledgement",                 // zlib/libpng License with Acknowledgement
      "ZPL-1.1",                              // Zope Public License 1.1
      "ZPL-2.0",                              // Zope Public License 2.0
      "ZPL-2.1",                              // Zope Public License 2.1
    ];
    $this->licenses_numbers = count($this->licenses);
    
    $this->main_copyright_owner = [
      'City of Roma',
      'City of Milano',
      'City of Napoli',
      'City of Torino',
      'City of Palermo',
      'City of Genova',
      'City of Bologna',
      'City of Firenze',
      'City of Bari',
      'City of Catania',
      'City of Venezia',
      'City of Verona',
      'City of Messina',
      'City of Padova',
      'City of Trieste',
      'City of Taranto',
      'City of Brescia',
      'City of Parma',
      'City of Prato',
      'City of Modena',
      'City of Reggio Calabria',
      'City of Reggio Emilia',
      'City of Perugia',
      'City of Livorno',
      'City of Ravenna',
      'City of Cagliari',
      'City of Foggia',
      'City of Rimini',
      'City of Salerno',
      'City of Ferrara',
      'City of Sassari',
      'City of Latina',
      'City of Giugliano in Campania',
      'City of Monza',
      'City of Siracusa',
      'City of Pescara',
      'City of Bergamo',
      'City of Forlì',
      'City of Trento',
      'City of Vicenza',
      'City of Terni',
      'City of Bolzano',
      'City of Novara',
      'City of Piacenza',
      'City of Ancona',
      'City of Andria',
      'City of Arezzo',
      'City of Udine',
      'City of Cesena',
      'City of Barletta',
    ];
    $this->main_copyright_owner_numbers = count($this->main_copyright_owner);

    $this->repo_owner = [
      'City of Roma',
      'City of Milano',
      'City of Napoli',
      'City of Torino',
      'City of Palermo',
      'City of Genova',
      'City of Bologna',
      'City of Firenze',
      'City of Bari',
      'City of Catania',
      'City of Venezia',
      'City of Verona',
      'City of Messina',
      'City of Padova',
      'City of Trieste',
      'City of Taranto',
      'City of Brescia',
      'City of Parma',
      'City of Prato',
      'City of Modena',
      'City of Reggio Calabria',
      'City of Reggio Emilia',
      'City of Perugia',
      'City of Livorno',
      'City of Ravenna',
      'City of Cagliari',
      'City of Foggia',
      'City of Rimini',
      'City of Salerno',
      'City of Ferrara',
      'City of Sassari',
      'City of Latina',
      'City of Giugliano in Campania',
      'City of Monza',
      'City of Siracusa',
      'City of Pescara',
      'City of Bergamo',
      'City of Forlì',
      'City of Trento',
      'City of Vicenza',
      'City of Terni',
      'City of Bolzano',
      'City of Novara',
      'City of Piacenza',
      'City of Ancona',
      'City of Andria',
      'City of Arezzo',
      'City of Udine',
      'City of Cesena',
      'City of Barletta',
    ];
    $this->repo_owner_numbers = count($this->repo_owner);

    $this->maintainance_type = ["internal", "contract", "community", "none"];
    $this->maintainance_type_numbers = count($this->maintainance_type);
    
    $this->technical_contacts = [];

    $this->platforms = [
      'web',
      'linux',
      'windows',
      'mac',
      'android',
      'ios'
    ];
    $this->platforms_numbers = count($this->platforms);

    $this->scope = ["it", "en", "es", "fr", "de", "sv", "sl", "nl", "cs"];
    $this->scope_numbers = count($this->scope);

    $this->ISO_639_3 = ["ita", "eng", "spa", "fra", "deu", "swe", "slv", "nld", "ces"];
    $this->ISO_639_3_numbers = count($this->ISO_639_3);

    $this->patype = [
      'city',
      'hospital',
      'police',
      'school',
      'university',
      'it-ag-turismo',
      'it-ag-lavoro',
      'it-ag-agricolo',
      'it-ag-formazione',
      'it-ag-fiscale',
      'it-ag-negoziale',
      'it-ag-erogagric',
      'it-ag-sanita',
      'it-ag-dirstudio',
      'it-altrilocali',
      'it-aci',
      'it-au-indip',
      'it-au-ato',
      'it-au-bacino',
      'it-au-portuale',
      'it-az-edilizia',
      'it-az-autonomo',
      'hospital',
      'it-az-servizi',
      'it-az-sanita',
      'it-camcom',
      'it-metro',
      'city',
      'it-montana',
      'it-co-bacino',
      'it-co-ricerca',
      'it-co-industria',
      'it-co-locali',
      'it-centrale',
      'it-provincia',
      'police',
      'it-regione',
      'it-afam',
      'school',
      'university',
    ];
    $this->patype_numbers = count($this->patype);

    $this->usedby = [
      'Comune di Roma',
      'Comune di Milano',
      'Comune di Napoli',
      'Comune di Torino',
      'Comune di Palermo',
      'Comune di Genova',
      'Comune di Bologna',
      'Comune di Firenze',
      'Comune di Bari',
      'Comune di Catania',
      'Comune di Venezia',
      'Comune di Verona',
      'Comune di Messina',
      'Comune di Padova',
      'Comune di Trieste',
      'Comune di Taranto',
      'Comune di Brescia',
      'Comune di Parma',
      'Comune di Prato',
      'Comune di Modena',
      'Comune di Reggio Calabria',
      'Comune di Reggio Emilia',
      'Comune di Perugia',
      'Comune di Livorno',
      'Comune di Ravenna',
      'Comune di Cagliari',
      'Comune di Foggia',
      'Comune di Rimini',
      'Comune di Salerno',
      'Comune di Ferrara',
      'Comune di Sassari',
      'Comune di Latina',
      'Comune di Giugliano in Campania',
      'Comune di Monza',
      'Comune di Siracusa',
      'Comune di Pescara',
      'Comune di Bergamo',
      'Comune di Forlì',
      'Comune di Trento',
      'Comune di Vicenza',
      'Comune di Terni',
      'Comune di Bolzano',
      'Comune di Novara',
      'Comune di Piacenza',
      'Comune di Ancona',
      'Comune di Andria',
      'Comune di Arezzo',
      'Comune di Udine',
      'Comune di Cesena',
      'Comune di Barletta',
    ];
    $this->usedby_numbers = count($this->usedby);

    $this->tags = [
      'rest',
      'fixed',
      'jewel',
      'head',
      'imaginary',
      'mellow',
      'seashore',
      'intend',
      'square',
      'decorate',
      'attach',
      'pie',
      'thank',
      'wanting',
      'oil',
      'crooked',
      'coil',
      'furry',
      'surround',
      'complete',
    ];
    $this->tags_numbers = count($this->tags);

    $this->category = [
      'it-mammoth',
      'it-giant',
      'it-spotty',
      'it-boundless',
      'fr-thoughtful',
      'fr-barbarous',
      'en-languid',
      'en-chunky',
      'en-dizzy',
      'de-unsightly',
      'de-sore',
      'en-fallacious',
    ];
    $this->category_numbers = count($this->category);

    $this->dependencies = [
      'Oracle',
      'MySQL',
      'Apache',
      'Varnish',
      'Docker',
      'Redis',
      'MS SQL',
      'nginx',
    ];
    $this->dependencies_numbers = count($this->dependencies);

    $this->dependencies_hardware = [
      'NFC Reader (chipset xxx)'
    ];
    $this->dependencies_hardware_numbers = count($this->dependencies);

    $this->metadata_repo = $this->readExampleMetadataRepo();
    $this->metadata_repo_numbers = count($this->metadata_repo);

    $this->development_status = [
      "concept",
      "development",
      "beta",
      "stable",
      "obsolete"
    ];
    $this->development_status_numbers = count($this->development_status);

    $this->software_type = [
      "standalone",
      "addon",
      "library",
      "configurationFiles"
    ];
    $this->software_type_numbers = count($this->software_type);

    $this->mime_types = [
      "audio/aac",
      "application/x-abiword",
      "application/octet-stream",
      "video/x-msvideo",
      "application/vnd.amazon.ebook",
      "application/octet-stream",
      "image/bmp",
      "application/x-bzip",
      "application/x-bzip2",
      "application/x-csh",
      "text/css",
      "text/csv",
      "application/msword",
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "application/vnd.ms-fontobject",
      "application/epub+zip",
      "application/ecmascript",
      "image/gif",
      "text/html",
      "image/x-icon",
      "text/calendar",
      "application/java-archive",
      "image/jpeg",
      "application/javascript",
      "application/json",
      "audio/midi",
      "audio/x-midi",
      "video/mpeg",
      "application/vnd.apple.installer+xml",
      "application/vnd.oasis.opendocument.presentation",
      "application/vnd.oasis.opendocument.spreadsheet",
      "application/vnd.oasis.opendocument.text",
      "audio/ogg",
      "video/ogg",
      "application/ogg",
      "font/otf",
      "image/png",
      "application/pdf",
      "application/vnd.ms-powerpoint",
      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
      "application/x-rar-compressed",
      "application/rtf",
      "application/x-sh",
      "image/svg+xml",
      "application/x-shockwave-flash",
      "application/x-tar",
      "image/tiff",
      "application/typescript",
      "font/ttf",
      "application/vnd.visio",
      "audio/wav",
      "audio/webm",
      "video/webm",
      "image/webp",
      "font/woff",
      "font/woff2",
      "application/xhtml+xml",
      "application/vnd.ms-excel",
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "application/xml",
      "application/vnd.mozilla.xul+xml",
      "application/zip",
      "video/3gpp",
      "audio/3gpp",
      "application/x-7z-compressed",
    ];
    $this->mime_types_numbers = count($this->mime_types);

    $this->ecosistemi = [
      'sanita',
      'welfare',
      'finanza-pubblica',
      'scuola',
      'istruzione-superiore-ricerca',
      'difesa-sicurezza-soccorso-legalita',
      'giustizia',
      'infrastruttura-logistica',
      'sviluppo-sostenibilita',
      'beni-culturali-turismo',
      'agricoltura',
      'italia-europa-mondo',
    ];
    $this->ecosistemi_numbers = count($this->ecosistemi);
  }

  public function generateDocuments($n = 100) {
    $documents = [];
    $this->descriptions = $this->getDocumentsDescription($n);
    
    // 1 January 2005 01:01:01
    $start = 1104541261;

    $now = (new DateTime())->getTimestamp();

    // 1 January 2025 01:01:01
    $end = 1735693261;

    for ($i=0; $i < $n; $i++) { 
      $name = $this->getRandomProjectName();
      $audience_countries = $this->getRandomScope();
      $audience_unsupported_countries = $this->getRandomScope();
      $intended_audience_unsupported_countries = [];
      foreach ($audience_unsupported_countries as $value) {
        if (!in_array($value, $audience_countries)) {
          $intended_audience_unsupported_countries[] = $value;
        }
      }

      $tags = $this->getRandomTags();
      $free_tags = $this->getRandomFreeTags();

      $documents[] = [
        "publiccode-yaml-version" => "http://w3id.org/publiccode/version/0.1",
        "name" => $name,
        "application-suite" => $this->getRandomApplicationSuite(),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$name.".git",
        "landing-url" => "https://example.com/italia/medusa",
        "is-based-on" => $this->getRandomIsBasedOn(),
        "software-version" => $this->getRandomVersion(),
        "release-date" => $this->getRandomDate($start, $now),
        "logo" => "img/logo.svg",
        "monochrome-logo" => "img/logo-mono.svg",
        "input-types" => $this->getRamdomMimeTypes(),
        "output-types" => $this->getRamdomMimeTypes(0, 3),
        "platforms" => $this->getRandomPlatforms(),
        "tags" => $tags,
        "free-tags" => $free_tags,
        "used-by" => $this->getRandomUsedBy(),
        "roadmap" => "https://example.com/italia/medusa/roadmap",
        "development-status" => $this->getRandomDevelopmentStatus(),
        "software-type" => $this->getRandomSoftwareType(),
        "software-is-related-to" => $this->getRandomSoftwareIsRelatedTo(),
        "intended-audience-only-for" => $this->getRandomPaType(),
        "intended-audience-countries" => $audience_countries,
        "intended-audience-unsupported-countries" => $intended_audience_unsupported_countries,
        "legal-license" => $this->getRandomLicense(),
        "legal-main-copyright-owner" => $this->getRandomMainCopyrightOwner(),
        "legal-repo-owner" => $this->getRandomMainCopyrightOwner(),
        "legal-authors-file" => "doc/AUTHORS.txt",
        "description" => [
          "ita" => $this->getRandomDescription($name, $i),
          "eng" => $this->getRandomDescription($name, $i),
        ],
        "dependencies-open" => $this->getRandomDependencies(),
        "dependencies-proprietary" => $this->getRandomDependencies(),
        "dependencies-hardware" => $this->getRandomDependenciesHardware(),
        "maintainance-contacts" => $this->generateRandomMaintenanceContact(),
        "maintainance-contractors" => $this->getRandomMaintenanceContractors(),
        "maintainance-type" => $this->getRandomMaintainanceType(),
        "localisation-localisation-ready" => boolval(rand(0,1)),
        "localisation-available-languages" => [],
        "it-conforme-accessibile" => boolval(rand(0,1)),
        "it-conforme-interoperabile" => boolval(rand(0,1)),
        "it-conforme-sicuro" => boolval(rand(0,1)),
        "it-conforme-privacy" => boolval(rand(0,1)),
        "it-spid" => boolval(rand(0,1)), 
        "it-cie" => boolval(rand(0,1)), 
        "it-anpr" => boolval(rand(0,1)), 
        "it-pagopa" => boolval(rand(0,1)), 
        "it-riuso-codice-ipa" => $this->generateRandomString(6, TRUE),
        "it-ecosistemi" => $this->getRandomEcosistemi(),
        "it-design-kit-seo"  => boolval(rand(0,1)), 
        "it-design-kit-ui"  => boolval(rand(0,1)), 
        "it-design-kit-web" => boolval(rand(0,1)),
        "it-design-kit-content" => boolval(rand(0,1)),
        "suggest-name" => explode(" ", $name),
        "metadata-repo" => $this->getRandomMetadataRepo(),
        "vitality-score" => rand(1, 100),
        "vitality-data-chart" => $this->getRandomVitalityDataChart(),
        "tags-related" => $this->getRandomItemFromArray($tags),
        "popular-tags" => $this->getRandomItemFromArray($tags),
        "share-tags" => $this->getRandomItemFromArray($tags),
        "related-software" => $this->getRandomRelatedSoftware(),
        "old-variant" => $this->getRandomOldVariant(),
      ];
    }

    return $documents;
  }

  public function getRandomLicense() {
    return $this->licenses[rand(0, $this->licenses_numbers - 1)];
  }

  public function getRandomMainCopyrightOwner() {
    return $this->main_copyright_owner[rand(0, $this->main_copyright_owner_numbers - 1)];
  }

  public function getRandomVersion() {
    $maj = rand(1,3);
    $min = rand(1,30);
    $build = rand(1,1000);

    return $maj . "." . $min . "." . $build;
  }

  public function getRandomDate($start, $end) {
    $timestamp = rand($start, $end);
    return date("Y-m-d", $timestamp);
  }

  public function getRandomVideoUrls() {
    $n = rand(0,3);
    $videos = [];
    for ($i=0; $i < $n; $i++) { 
      $videos[] = 'https://youtube.com/' . $this->generateRandomString(8);
    }

    return $videos;
  }

  public function getRandomPlatforms() {
    $n = rand(1, $this->platforms_numbers);
    $platforms = [];

    for ($i=0; $i < $n;) { 
      $current = rand(0, $this->platforms_numbers - 1);
      if(!in_array($this->platforms[$current], $platforms)) {
        $platforms[] = $this->platforms[$current];
        $i++;
      }
    }

    return $platforms;
  }

  public function getRandomScope() {
    $n = rand(1, $this->scope_numbers);
    $scope = [];

    for ($i=0; $i < $n;) { 
      $current = rand(0, $this->scope_numbers - 1);
      if(!in_array($this->scope[$current], $scope)) {
        $scope[] = $this->scope[$current];
        $i++;
      }
    }

    return $scope;
  }

  public function getRandomPaType() {
    $n = rand(1, 10);
    $patype = [];

    for ($i=0; $i < $n;) { 
      $current = rand(0, $this->patype_numbers - 1);
      if(!in_array($this->patype[$current], $patype)) {
        $patype[] = $this->patype[$current];
        $i++;
      }
    }

    return $patype;
  }

  public function getRandomUsedBy() {
    $n = rand(1, 8);
    $usedby = [];

    for ($i=0; $i < $n;) { 
      $current = rand(0, $this->usedby_numbers - 1);
      if(!in_array($this->usedby[$current], $usedby)) {
        $usedby[] = $this->usedby[$current];
        $i++;
      }
    }

    return $usedby;
  }

  public function getRandomTags() {
    $n = rand(1, 15);
    $tags = [];

    for ($i=0; $i < $n;) { 
      $current = rand(0, $this->tags_numbers - 1);
      if(!in_array($this->tags[$current], $tags)) {
        $tags[] = $this->tags[$current];
        $i++;
      }
    }

    return $tags;
  }

  public function getRandomFreeTags() {
    $free_tags = [];
    // numero di tags
    $n_tags = rand(0, 9);
    $n_langs = rand(0, 5);
    if($n_tags > 0) {
      for ($i=0; $i < $n_langs; $i++) { 
        $free_tags[$this->getRandomISO6393()] = $this->generateRandomFreeTags($n_tags);
      }
    }

    return $free_tags;
  }

  public function generateRandomFreeTags($n = 3) {
    $free_tags = [];

    for ($i=0; $i < $n; $i++) { 
      $free_tags[] = $this->getRandomPhrase(2, 3, "-");
    }

    return $free_tags;
  }

  public function getRandomDependencies() {
    $n = rand(0, 5);
    $dependencies = [];

    for ($i=0; $i < $n;) {
      $current = rand(0, $this->dependencies_numbers - 1);
      if(!in_array($this->dependencies[$current], $dependencies)) {
        if(($i % 2) == 0){
          $dependencies[] = [
            'name' => $this->dependencies[$current],
            'version' => $this->getRandomVersion(),
            'optional' => boolval(rand(0,1))
          ];
        }
        else {
          $dependencies[] = [
            'name' => $this->dependencies[$current],
            'version-min' => $this->getRandomVersion(),
            'version-max' => $this->getRandomVersion(),
            'optional' => boolval(rand(0,1))
          ];
        }
        $i++;
      }
    }

    return $dependencies;
  }

  public function getRandomDependenciesHardware() {
    $n = rand(0, 1);
    $dependencies_hardware = [];

    if ($n == 1) {
      return [
        'name' => $this->dependencies_hardware,
        'optional' => boolval(rand(0,1)),
      ];
    }

    return $dependencies_hardware;
  }

  public function getRandomMaintainanceType() {
    return $this->maintainance_type[rand(0, $this->maintainance_type_numbers -1)];
  }

  public function getRandomMaintenanceContractors() {
    $maintainance_contractors = [];
    $n = rand(1,3);
    for ($i=0; $i < $n; $i++) { 
      $maintainance_contractors = $this->generateRandomMaintenanceContractor();
    }

    return $maintainance_contractors;
  }

  public function getRandomProjectName() {
    $n = rand(1, 4);
    $name = "";

    for ($i=0; $i < $n; $i++) { 
      $name .= ucfirst(strtolower($this->generateRandomString(rand(4, 10), TRUE))) . " ";
    }

    return trim($name);
  }

  public function getRandomApplicationSuite() {
    $n = rand(1, 4);
    $application_suite = "";

    for ($i=0; $i < $n; $i++) { 
      $application_suite .= ucfirst(strtolower($this->generateRandomString(rand(4, 10), TRUE)));
    }

    return $application_suite;
  }

  public function getRandomIsBasedOn() {
    $n = rand(1, 4);
    $is_based_on = [];

    for ($i=0; $i < $n; $i++) {
      $is_based_on[] = "https://github.com/italia/otello.git";
    }

    return $is_based_on;
  }

  public function getRandomDevelopmentStatus() {
    return $this->development_status[rand(0, $this->development_status_numbers -1)];
  }

  public function getRandomSoftwareType() {
    return $this->software_type[rand(0, $this->software_type_numbers - 1)];
  }

  public function getRandomSoftwareIsRelatedTo() {
    $n = rand(0, 4);
    $software_is_related_to = [];

    for ($i=0; $i < $n ; $i++) { 
      $software_is_related_to[] = "https://example.com/italia/software".$i.".git";
    }

    return $software_is_related_to;
  }

  public function getRandomDescription($name, $i_description) {
    $screenshots = [];
    $n = rand(1, 8);
    for ($i=0; $i < $n; $i++) { 
      $screenshots[] = "img/sshot".($i+1).".jpg";
    }
    return [
      "localised-name" => $name,
      "generic-name" => $this->getRandomPhrase(1, 3),
      "short-description" => substr($this->descriptions[$i_description*2], 0, rand(100, 150)),
      "long-description" => $this->descriptions[$i_description*2],
      "documentation" => "https://read.the.documentation/medusa/v1.0",
      "api-documentation" => "https://read.the.api-documentation/medusa/v1.0",
      "feature-list" => $this->getRandomFeatureList(),
      "screenshots" => $screenshots,
      "videos" => $this->getRandomVideoUrls(),
      "awards" => $this->getRandomAwardsList(),
    ];
  }

  public function getRandomISO6393() {
    return $this->ISO_639_3[rand(0, ($this->ISO_639_3_numbers - 1))];
  }

  public function getDocumentsDescription($n = 100) {
    
    $descriptions = [];
    $numbers = $n * 2;
    $retry = 3;

    while($numbers > 0 && $retry > 0) {
      $curl = curl_init();
      curl_setopt_array($curl, array(
          CURLOPT_RETURNTRANSFER => 1,
          CURLOPT_URL => 'https://baconipsum.com/api/?type=all-meat&paras='.$numbers.'&format=json'
      ));

      try {
        $resp = curl_exec($curl);
        $resp = json_decode($resp);
        $descriptions = array_merge($descriptions, $resp);

        $retry = 3;
        $numbers = $numbers - count($resp);
      }
      catch(Exeption $e) {
        $retry--;
      }
    }

    return $descriptions;
  }

  public function getRandomMetadataRepo() {
    return $this->metadata_repo[rand(0, $this->metadata_repo_numbers - 1)];
  }

  public function getRandomVitalityDataChart() {
    $months = rand(0, 12);
    $vitality_data_chart = [];
    for ($i=0; $i < $months ; $i++) { 
      $vitality_data_chart[] = rand(1, 10);
    }

    return $vitality_data_chart;
  }

  public function getRandomItemFromArray($all_tags = []) {
    if(empty($all_tags)){
      return [];
    }
    
    $c = count($all_tags);
    $n = rand(1, $c);
    
    $keys = array_rand($all_tags, $n);
    if(!is_array($keys)){
      return [
        $all_tags[$keys]
      ];
    }

    $items = [];    
    foreach ($keys as $key) {
      $items[] = $all_tags[$key];
    }
    
    return $items;
  }

  public function getRandomRelatedSoftware() {
    $n = rand(0, 3);
    $items = [];

    for ($i=0; $i < $n; $i++) { 
      $items[] = $this->generateRandomRelatedSoftware();
    }

    return $items;
  }

  public function getRandomOldVariant() {
    $n = rand(0, 3);
    $items = [];

    for ($i=0; $i < $n; $i++) { 
      $items[] = $this->generateRandomOldVariant();
    }

    return $items;
  }

  public function getRamdomMimeTypes($min = 0, $max = 5) {
    $n = rand($min, $max);
    $mime_types = [];

    for ($i=0; $i < $n; $i++) { 
      $mime_types[] = $this->mime_types[rand(0, $this->mime_types_numbers - 1)];
    }

    return $mime_types;
  }

  public function getRandomEcosistemi($min = 0, $max = 5) {
    $n = rand($min, $max);
    $ecosistemi = [];

    for ($i=0; $i < $n; $i++) { 
      $ecosistemi[] = $this->ecosistemi[rand(0, $this->ecosistemi_numbers - 1)];
    }

    return $ecosistemi;
  }

  private function generateRandomString($length = 10, $only_letters = FALSE) {
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    if ($only_letters) {
      $characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';      
    }
    $charactersLength = strlen($characters);
    $randomString = '';

    for ($i = 0; $i < $length; $i++) {
        $randomString .= $characters[rand(0, $charactersLength - 1)];
    }

    return $randomString;
  }

  private function getRandomAwardsList() {
    $n = rand(0, 3);
    $awards_list = [];

    for ($i=0; $i < $n; $i++) { 
      $awards_list[] = $this->getRandomPhrase(2, 4);
    }

    return $awards_list;
  }

  private function getRandomFeatureList() {
    $n = rand(1, 5);
    $feature_list = [];

    for ($i=0; $i < $n; $i++) { 
      $feature_list[] = $this->getRandomPhrase();
    }

    return $feature_list;
  }

  private function getRandomPhrase($min = 3, $max = 7, $blank = " ") {
    $n = rand($min, $max);
    $feature = "";

    for ($i=0; $i < $n; $i++) { 
      $feature .= strtolower($this->generateRandomString(rand(2, 12), TRUE)) . $blank;
    }

    return trim($feature, $blank);
  }

  private function generateRandomMaintenanceContact() {
    $name = strtolower($this->generateRandomString(rand(6, 15), TRUE));
    $surname = strtolower($this->generateRandomString(rand(6, 15), TRUE));
    $affiliation = $this->generateRandomAffiliation();
    return [
      "name" => ucfirst($name) . " " . ucfirst($surname),
      "email" => $name.".".$surname."@example.com",
      "phone" => "123456789",
      "affiliation" => $this->generateRandomAffiliation(),
    ];
  }

  private function generateRandomMaintenanceContractor() {
    $now = (new DateTime())->getTimestamp();
    $end = 1735693261;

    return [
      "name" => ucfirst(strtolower($this->generateRandomString(rand(6, 15), TRUE))) . " " . ucfirst(strtolower($this->generateRandomString(rand(6, 15), TRUE))) . " S.p.A.",
      "website" => "https://www.companywebsite.com",
      "until" => $this->getRandomDate($now, $end),
    ];
  }

  private function generateRandomAffiliation() {
    $n = rand(2, 4);
    $affiliation = '';
    for ($i=0; $i < $n; $i++) { 
      $affiliation .= ucfirst(strtolower($this->generateRandomString(rand(6, 15), TRUE))) . " ";
    }

    return trim($affiliation);
  }

  private function readExampleMetadataRepo() {
    $metadata = [];
    $files = [
      'metadata-repo-bitbucket.json',
      'metadata-repo-github.json',
      'metadata-repo-gitlab.json',
    ];

    foreach ($files as $file) {
      $json = file_get_contents($file);
      $metadata[] = json_decode($json);
    }

    return $metadata;
  }

  private function generateRandomRelatedSoftware() {
    return [
      "name" => $this->generateRandomString(rand(6, 15), TRUE),
      "image" => "img/screenshot.jpg",
      "eng" => [
        "localised-name" => $this->getRandomPhrase(1, 3),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$this->generateRandomString(rand(5, 10), TRUE).".git",
      ],
      "ita" => [
        "localised-name" => $this->getRandomPhrase(1, 3),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$this->generateRandomString(rand(5, 10), TRUE).".git",
      ],
    ];
  }

  private function generateRandomOldVariant() {
    return [
      "eng" => [
        "localised-name" => $this->getRandomPhrase(1, 3),
        "feature-list" => $this->getRandomFeatureList(),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$this->generateRandomString(rand(5, 10), TRUE).".git",
      ],
      "ita" => [
        "localised-name" => $this->getRandomPhrase(1, 3),
        "feature-list" => $this->getRandomFeatureList(),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$this->generateRandomString(rand(5, 10), TRUE).".git",
      ],
    ];
  }

}
