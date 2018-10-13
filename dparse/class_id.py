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

import re
from text import ascii_only


class ClassIdList(object):
    def __init__(self):
        self._entities = dict()     # Key -> (int) class_id,  Value-> (ClassId)

    def __getitem__(self, item):
        return self._entities[item]  # delegate to li.__getitem__

    def __iter__(self):
        for entry in self._entities:
            yield entry

    def __len__(self):
        return len(self._entities)

    def add(self, class_id):
        assert isinstance(class_id, ClassId), 'Invalid type'
        assert class_id.cid not in self._entities, 'Entity already exists: {}'.format(class_id)
        self._entities[class_id.cid] = class_id
        return self


class ClassId(object):
    def __init__(self):
        self.cid = None                 # Class Id
        self.name = None                # Title
        self.section = None             # Document section
        self.heading_para_no = None     # Paragraph number of section with heaself.ding

    def __str__(self):
        return 'ClassId: {}: {}'.format(self.cid, self.name)

    @staticmethod
    def create(number, paragraph):
        section = ClassId()

        section.paragraph_number = number
        section.style_name = paragraph.style.name

        try:
            assert 'heading ' in section.style_name.lower(), 'Heading style not found'
            if isinstance(paragraph.text, list):
                heading_txt = " ".join(paragraph.text)
            else:
                heading_txt = paragraph.text

            # scrub non-ascii
            heading_txt = ascii_only(heading_txt)

            # Split into section number and section title
            headings = re.findall(r"[^\s']+", heading_txt)

            assert len(headings) >= 2, \
                "Foreign Heading string format: '{}'".format(paragraph.text)

            section.section_number = headings[0]
            section.title = ' '.join(headings[1:])
            section.section_points = [pt for pt in section.section_number.split('.')]

            return section

        except Exception as _e:
            raise
