<?php

class generatorRandomDocuments {

  protected $licenses;
  protected $licenses_numbers;

  protected $main_copyright_owner;
  protected $main_copyright_owner_numbers;

  protected $repo_owner;
  protected $repo_owner_numbers;

  protected $maintenance_type;
  protected $maintenance_type_numbers;

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

  protected $codiceIPA;
  protected $codiceIPA_numbers;

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

    $this->maintenance_type = ["internal", "contract", "community", "none"];
    $this->maintenance_type_numbers = count($this->maintenance_type);

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
      // International tags.
    	"3dgraphics",    // application for viewing, creating, or processing 3-d graphics
    	"accessibility", // accessibility
    	"accounting",    // accounting software
    	"amusement",     // a simple amusement
    	"archiving",     // a tool to archive/backup data
    	"art",           // software to teach arts
    	"artificial-intelligence", // artificial intelligence software
    	"backend",                 // software not meant for end users
    	"calculator",              // a calculator
    	"calendar",                // calendar application
    	"chat",                    // a chat client
    	"classroom-management",    // classroom management software
    	"clock",                   // a clock application/applet
    	"content-management",      // a content management system (CMS)
    	"compression",             // a tool to manage compressed data/archives
    	"construction",            //
    	"contact-management",      // e.g. an address book
    	"database",                // application to manage a database
    	"debugger",                // a tool to debug applications
    	"dictionary",              // a dictionary
    	"documentation",           // help or documentation
    	"electronics",             // electronics software, e.g. a circuit designer
    	"email",                   // email application
    	"emulator",                // emulator of another platform, such as a dos emulator
    	"engineering",             // engineering software, e.g. cad programs
    	"file-manager",            // a file manager
    	"file-transfer",           // tools like ftp or p2p programs
    	"finance",                 // application to manage your finance
    	"flowchart",               // a flowchart application
    	"gui-designer",            // a gui designer application
    	"identity",                // identity management
    	"instant-messaging",       // an instant messaging client
    	"library",                 // a library software
    	"medical",                 // medical software
    	"monitor",                 // monitor application/applet that monitors some resource or activity
    	"museum",                  // museum software
    	"music",                   // musical software
    	"news",                    // software to manage and publish news
    	"ocr",                     // optical character recognition application
    	"parallel-computing",      // parallel computing software
    	"photography",             // camera tools, etc.
    	"presentation",            // presentation software
    	"printing",                // a tool to manage printers
    	"procurement",             // software for managing procurement
    	"project-management",      // project management application
    	"publishing",              // desktop publishing applications and color management tools
    	"raster-graphics",         // application for viewing, creating, or processing raster (bitmap) graphics
    	"remote-access",           // a tool to remotely manage your pc
    	"revision-control",        // applications like git or subversion
    	"robotics",                // robotics software
    	"scanning",                // tool to scan a file/text
    	"security",                // a security tool
    	"sports",                  // sports software
    	"spreadsheet",             // a spreadsheet
    	"telephony",               // telephony via pc
    	"terminal-emulator",       // a terminal emulator application
    	"texteditor",              // a text editor
    	"texttools",               // a text tool utility
    	"translation",             // a translation tool
    	"vector-graphics",         // application for viewing, creating, or processing vector graphics
    	"video-conference",        // video conference software
    	"viewer",                  // tool to view e.g. a graphic or pdf file
    	"web-browser",             // a web browser
    	"whistleblowing",          // software for whistleblowing / anticorruption
    	"word-processor",          // a word processor
    	"wordprocessor",           // a word processor
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

    $this->codiceIPALabel = [
      '054' => "Azienda Unita' Sanitaria Locale Umbria 1",
      '055' => "Azienda Unita' Sanitaria Locale Umbria 2",
      '056' => "Azienda Unita' Sanitaria Locale Viterbo",
      '058' => "Azienda Sanitaria Locale Roma 5",
      '080' => "Azienda Sanitaria Provinciale N. 5 di Reggio Calabria",
      '092' => "Azienda Ospedaliera Brotzu",
      '093' => "Azienda Assistenza Sanitaria N.5 Friuli Occidentale",
      '102' => "Azienda Sanitaria Provinciale di Vibo Valentia",
      '1CMVL' => "Comunita' Montana Valle del Liri Zona XV",
      'a13_f952' => "Azienda Sanitaria Locale NO",
      'a1_025' => "Azienda ULSS n. 1 Dolomiti",
      'a1_na' => "Azienda Sanitaria Locale Napoli 1 Centro",
      'a1_ss' => "Azienda per la tutela della salute",
      'a4_te' => "Azienda Sanitaria Locale N. 4 di Teramo",
      'A690_bpe' => "Comune di Bascape'",
      'a9_tv' => "Azienda ULSS n. 2 Marca Trevigiana",
      'aaba_001' => "Accademia Albertina di Belle Arti di Torino",
      'aaca' => "Associazione Ambito Cuneese Ambiente",
      'aaci_' => "Autorita di Ambito Calore Irpino",
      'aacsl' => "Ato Ambiente Cl2 Spa in Liquidazione",
      'aacsl1' => "Ato Ambiente Cl1 Spa in Liquidazione",
      'aacst' => "Azienda Autonoma di Cura Soggiorno e Turismo Castellammare Di Stabia",
      'aacstip' => "Azienda Autonoma di Cura Soggiorno e Turismo Delle Isole di Ischia e di Procida",
      'aacstpoz' => "Azienda Autonoma di Cura Soggiorno e Turismo di Pozzuoli",
      'aacst_' => "Azienda Autonoma di Cura Soggiorno e Turismo di Vico Equense",
      'aacst_na' => "Azienda Autonoma di Cura Soggiorno e Turismo di Capri",
      'aacst_po' => "Azienda Autonoma di Cura Soggiorno e Turismo di Pompei",
      'aadtaci' => "Azienda Autonoma Delle Terme di Acireale",
      'aafc' => "Azienda Assistenza Sanitaria 3 Alto Friuli Collinare Medio Friuli",
      'aaill_' => "Arlas - Agenzia per Il Lavoro e L'Istruzione",
      'aalss' => "Ales Arte Lavoro e Servizi Spa",
      'aamsc' => "Casa di Riposo Sgobba",
      'aao_al' => "Appennino Aleramico Obertengo di Ponzone",
      'aapl_001' => "Agenzia di Accoglienza e Promozione Turistica del Territorio della Provincia di Torino",
      'aaric' => "Aric Agenzia Regionale di Informatica e Committenza",
      'aarspa' => "Alto Adige Riscossioni Spa",
      'aas' => "Azienda Assistenza Sanitaria N.2 Bassa Friulana Isontina",
      'aas1ts' => "Azienda Sanitaria Universitaria Integrata di Trieste",
      'aasct_na' => "Azienda Autonoma Soggiorno Cura e Turismo di Napoli",
      'aasetp' => "Azienda Autonoma Soggiorno e Turismo Paestum",
      'aasmc' => "Associazione Arena Sferisterio",
      'aasnq' => "Azienda Sanitaria Universitaria Integrata di Udine",
      'aasp' => "Asea Azienda Speciale",
      'aaspa' => "Anconambiente S.P.A.",
      'aasssa' => "Azienda Autonoma di Soggiorno di Sorrento-Sant'Agnello",
      'aast' => "Azienda Autonoma di Soggiorno e Turismo Termoli",
      'aastam' => "Azienda Autonoma Soggiorno e Turismo Amalfi",
      'aastc' => "Azienda Autonoma Soggiorno e Turismo Cava Dei Tirreni",
      'aastma' => "Azienda Autonoma Soggiorno e Turismo di Maiori",
      'aastp' => "Azienda Autonoma Soggiorno e Turismo di Positano",
      'aastr' => "Azienda Autonoma Soggiorno e Turismo di Ravello",
      'aastsa' => "Aziebda Autonoma di Soggiorno e Turismo Salerno",
      'aatgr_to' => "Associazione D'Ambito Torinese Governo Rifiuti",
      'aatm' => "Atm Azienda Trasporti Messina",
      'aato1mnp' => "Assemblea di Ambito Territoriale Ottimale n.1 Marche Nord Pesaro e Urbino",
      'aato2mca' => "Assemblea di Ambito territoriale ottimale n. 2 Marche centro - Ancona",
      'aatolv' => "Consiglio di Bacino Laguna di Venezia",
      'aaton_' => "Assemblea di Ambito territoriale ottimale n. 3 Marche Centro Macerata",
      'aatopa' => "Autorita' Ambito Territoriale Ottimale 1 Palermo in Liquidazione",
      'aatoq' => "Autorita' di Ambito Territoriale Ottimale N. 4 Marche Centro Sud",
      'ab' => "Asuc di Bedollo",
      'abab_037' => "Accademia di Belle Arti di Bologna",
      'abab_072' => "Accademia di Belle Arti di Bari",
      'abac' => "Accademia di Belle Arti di Catania",
      'abac_045' => "Accademia di Belle Arti di Carrara",
      'abac_079' => "Accademia di Belle Arti di Catanzaro",
      'abaf_048' => "Accademia di Belle Arti di Firenze",
      'abaf_060' => "Accademia di Belle Arti di Frosinone",
      'abaf_071' => "Accademia di Belle Arti di Foggia",
      'abala_aq' => "Accademia di Belle Arti L'Aquila",
      'abam_043' => "Accademia di Belle Arti di Macerata",
      'aban' => "Accademia Belle Arti di Napoli",
      'abap_082' => "Accademia di Belle Arti di Palermo",
      'abar_' => "Accademia di Belle Arti di Roma",
      'abau' => "Accademia di Belle Arti di Urbino",
      'abav_027' => "Accademia di Belle Arti di Venezia",
      'aba_090' => "Accademia di Belle Arti Mario Sironi",
      'aba_80' => "Accademia di Belle Arti Reggio Calabria",
      'abbru' => "Asuc di Brusago",
      'abb_' => "Autorita' di Bacino della Basilicata",
      'abd' => "Amministrazione Beni Demaniali",
      'abdft_cb' => "Autorita' di Bacino Interregionale dei fiumi Trigno, Biferno e Minori, Saccione e Fortore",
      'abfa' => "Autorita' di Bacino del Fiume Arno",
      'abfitlpb' => "Autorita' di Bacino dei Fiumi Isonzo, Tagliamento, Livenza, Piave, Brenta-Bacchiglione",
      'abflgv' => "Autorita' di Bacino dei Fiumi Liri Garigliano e Volturno",
      'abft' => "Autorita di Bacino del Fiume Tronto",
      'abftrm' => "Autorita' di Bacino del Fiume Tevere",
      'abigs' => "Amministrazione Separata  Beni Usi Civici Cengles",
      'ablcpg' => "Autorita' di Bacino Lacuale Ceresio Piano e Ghirla",
      'abllmcmv' => "Autorita' di Bacino Lacuale Dei Laghi Maggiore Comabbio Monate Varese",
      'abl_097' => "Autorita' di Bacino del Lario e Dei Laghi Minori",
      'abnfa_tn' => "Autorita' di Bacino del Fiume Adige",
      'abpfs' => "Autorita' di Bacino Pilota del Fiume Serchio",
      'abrc' => "Autorita di Bacino Regionale della Campania Centrale",
      'abrcs' => "Autorita di Bacino Regionale di Campania Sud Ed Interregionale per Il Bacino Idrografico del Fiume Sele",
      'abspa' => "Abbanoa S.P.A.",
      'absrl' => "Acque Bresciane Srl",
      'abs_bg' => "Ateneo Bergamo Spa",
      'abucfu' => "Amministrazione Beni Uso Civico Frazione di Fucine",
      'abucgra' => "Amministrazione Separata Dei Beni di Uso Civico Graun",
    ];
    $this->codiceIPA = [
      '054',
      '055',
      '056',
      '058',
      '080',
      '092',
      '093',
      '102',
      '1CMVL',
      'a13_f952',
      'a1_025',
      'a1_na',
      'a1_ss',
      'a4_te',
      'A690_bpe',
      'a9_tv',
      'aaba_001',
      'aaca',
      'aaci_',
      'aacsl',
      'aacsl1',
      'aacst',
      'aacstip',
      'aacstpoz',
      'aacst_',
      'aacst_na',
      'aacst_po',
      'aadtaci',
      'aafc',
      'aaill_',
      'aalss',
      'aamsc',
      'aao_al',
      'aapl_001',
      'aaric',
      'aarspa',
      'aas',
      'aas1ts',
      'aasct_na',
      'aasetp',
      'aasmc',
      'aasnq',
      'aasp',
      'aaspa',
      'aasssa',
      'aast',
      'aastam',
      'aastc',
      'aastma',
      'aastp',
      'aastr',
      'aastsa',
      'aatgr_to',
      'aatm',
      'aato1mnp',
      'aato2mca',
      'aatolv',
      'aaton_',
      'aatopa',
      'aatoq',
      'ab',
      'abab_037',
      'abab_072',
      'abac',
      'abac_045',
      'abac_079',
      'abaf_048',
      'abaf_060',
      'abaf_071',
      'abala_aq',
      'abam_043',
      'aban',
      'abap_082',
      'abar_',
      'abau',
      'abav_027',
      'aba_090',
      'aba_80',
      'abbru',
      'abb_',
      'abd',
      'abdft_cb',
      'abfa',
      'abfitlpb',
      'abflgv',
      'abft',
      'abftrm',
      'abigs',
      'ablcpg',
      'abllmcmv',
      'abl_097',
      'abnfa_tn',
      'abpfs',
      'abrc',
      'abrcs',
      'abspa',
      'absrl',
      'abs_bg',
      'abucfu',
      'abucgra',
    ];
    $this->codiceIPA_numbers = count($this->codiceIPA);
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
      $codiceIPA = $this->getRandomCodiceIPA();

      $documents[] = [
        "fileRawURL" => "https://example.com/italia/medusa/publiccode.yml",
        "publiccode-yaml-version" => "http://w3id.org/publiccode/version/0.1",
        "name" => $name,
        "applicationSuite" => $this->getRandomApplicationSuite(),
        "url" => "https://example.com/".$this->generateRandomString(rand(5, 10), TRUE)."/".$name.".git",
        "landingURL" => "https://example.com/italia/medusa",
        "isBasedOn" => $this->getRandomIsBasedOn(),
        "softwareVersion" => $this->getRandomVersion(),
        "releaseDate" => $this->getRandomDate($start, $now),
        "logo" => "img/logo.svg",
        "monochromeLogo" => "img/logo-mono.svg",
        "inputTypes" => $this->getRamdomMimeTypes(),
        "outputTypes" => $this->getRamdomMimeTypes(0, 3),
        "platforms" => $this->getRandomPlatforms(),
        "tags" => $tags,
        "usedBy" => $this->getRandomUsedBy(),
        "roadmap" => "https://example.com/italia/medusa/roadmap",
        "developmentStatus" => $this->getRandomDevelopmentStatus(),
        "softwareType" => $this->getRandomSoftwareType(),
        "intendedAudience-onlyFor" => $this->getRandomPaType(),
        "intendedAudience-countries" => $audience_countries,
        "intendedAudience-unsupportedCountries" => $intended_audience_unsupported_countries,
        "legal-license" => $this->getRandomLicense(),
        "legal-mainCopyrightOwner" => $this->getRandomMainCopyrightOwner(),
        "legal-repoOwner" => $this->getRandomMainCopyrightOwner(),
        "legal-authorsFile" => "doc/AUTHORS.txt",
        "description" => [
          "ita" => $this->getRandomDescription($name, $i),
          "eng" => $this->getRandomDescription($name, $i),
        ],
        "dependsOn-open" => $this->getRandomDependencies(),
        "dependsOn-proprietary" => $this->getRandomDependencies(),
        "dependsOn-hardware" => $this->getRandomDependenciesHardware(),
        "maintenance-contacts" => $this->generateRandomMaintenanceContact(),
        "maintenance-contractors" => $this->getRandomMaintenanceContractors(),
        "maintenance-type" => $this->getRandommaintenanceType(),
        "localisation-localisationReady" => boolval(rand(0,1)),
        "localisation-availableLanguages" => [],
        "it-conforme-accessibile" => boolval(rand(0,1)),
        "it-conforme-interoperabile" => boolval(rand(0,1)),
        "it-conforme-sicuro" => boolval(rand(0,1)),
        "it-conforme-privacy" => boolval(rand(0,1)),
        "it-spid" => boolval(rand(0,1)),
        "it-cie" => boolval(rand(0,1)),
        "it-anpr" => boolval(rand(0,1)),
        "it-pagopa" => boolval(rand(0,1)),
        "it-riuso-codiceIPA" => $codiceIPA,
        "it-riuso-codiceIPA-label" => ($codiceIPA == NULL) ? $codiceIPA : $this->codiceIPALabel[$codiceIPA],
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
        "old-feature-list" => [
          "ita" => $this->getRandomFeatureList(),
          "eng" => $this->getRandomFeatureList(),
        ],
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

  public function getRandomCodiceIPA() {
    $codiceIPA = NULL;
    if (rand(0,1) == 1) {
      $codiceIPA = $this->codiceIPA[rand(0, ($this->codiceIPA_numbers - 1))];
    }

    return $codiceIPA;
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

  public function getRandommaintenanceType() {
    return $this->maintenance_type[rand(0, $this->maintenance_type_numbers -1)];
  }

  public function getRandomMaintenanceContractors() {
    $maintenance_contractors = [];
    $n = rand(1,3);
    for ($i=0; $i < $n; $i++) {
      $maintenance_contractors = $this->generateRandomMaintenanceContractor();
    }

    return $maintenance_contractors;
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

  public function getRandomDescription($name, $i_description) {
    $screenshots = [];
    $n = rand(1, 8);
    for ($i=0; $i < $n; $i++) {
      $screenshots[] = "img/sshot".($i+1).".jpg";
    }
    return [
      "localisedName" => $name,
      "genericName" => $this->getRandomPhrase(1, 3),
      "shortDescription" => substr($this->descriptions[$i_description*2], 0, rand(100, 150)),
      "longDescription" => $this->descriptions[$i_description*2],
      "documentation" => "https://read.the.documentation/medusa/v1.0",
      "apiDocumentation" => "https://read.the.api-documentation/medusa/v1.0",
      "featureList" => $this->getRandomFeatureList(),
      "freeTags" => $this->generateRandomFreeTags(rand(0, 9)),
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
    $old_variant = [
      "name" => $this->getRandomPhrase($min = 1, $max = 4),
      "vitality-score" => rand(1, 99),
      "legal-repo-owner" => $this->getRandomMainCopyrightOwner(),
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

    if (rand(0,1) == 0) {
      unset($old_variant["eng"]["localised-name"]);
      unset($old_variant["ita"]["localised-name"]);
    }

    return $old_variant;
  }

}
