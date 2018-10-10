#!/usr/bin/env python
#
# Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import types


def ascii_only(input_text):
    """
    Map Word Text to best ASCII equivalents
    See the 'findunicode.py' program for how to search for these

    :param input_text: (str) input text
    :return: (str) Ascii only text
    """
    replacements = {
        160  : u'-',
        174  : u'r',        # Registered sign
        176  : u"degree-",  # Degree sign
        177  : u"+/-",
        181  : u"u",        # Micro
        189  : u"1/2",
        215  : u'*',
        224  : u"`a",
        946  : u'B',        # Beta
        956  : u'v',
        969  : u'w',
        8211 : u'-',
        8217 : u"'",
        8220 : u"``",
        8221 : u"''",
        8230 : u"...",
        8722 : u'-',
        8804 : u'<=',
        61664: u'->',
        8805 : u'>=',
        8226 : u'o',        # Bullet
    }
    if isinstance(input_text, types.GeneratorType):
        text = ''.join(i for i in input_text)
    else:
        text = input_text

    return ''.join(
            [u'' if len(i) == 0
             else replacements[ord(i)] if ord(i) >= 128 and ord(i) in replacements
             else i if ord(i) < 128
             else u' ' for i in text])
