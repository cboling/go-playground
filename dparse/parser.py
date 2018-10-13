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
#
from __future__ import (
    absolute_import, division, print_function, unicode_literals
)

from docx import Document

from section import SectionHeading, SectionList
from class_id import ClassIdList, ClassId
from tables import TableList, Table
from text import ascii_only, camelcase


MESectionStartLabel = "9.1.1"   # First to read
MESectionEndLabel = "10"        # MEs not found after here
MEClassSection = "11.2.4"       # Class IDs


def parse_args():
    return {
        'itu': 'T-REC-G.988-201711-I!!MSW-E.docx',
        'sections': 'G.988.Sections.json',
        'preparsed': 'G.988.PreCompiled.json',
        'output': 'G.988.Parsed.json',
        'me-start': MESectionStartLabel,
        'me-end': MESectionEndLabel,
        'me-class-section': MEClassSection,
    }


class Main(object):
    """ Main program """
    def __init__(self):
        self.args = parse_args()
        self.paragraphs = None
        self.sections = None
        self.class_ids = None
        self.body = None

    def load_itu_document(self):
        return Document(self.args['itu'])

    def start(self):
        print("Loading ITU Document '{}' and parsed Section Header file '{}'".
              format(self.args['itu'], self.args['sections']))

        document = self.load_itu_document()
        self.sections = SectionList()
        self.sections.load(self.args['preparsed'])

        self.paragraphs = document.paragraphs
        # doc_sections = document.sections
        # styles = document.styles
        # self.body = document.element.body

        print('Extracting ME Class ID values')

        self.class_ids = ClassIdList.parse_sections(self.sections,
                                                    self.args['me-class-section'])

        print('Found {} ME Class ID entries. {} have sections associated to them'.
              format(len(self.class_ids), len([c for c in self.class_ids.values()
                                               if c.section is not None])))

        print('Managed Entities without Sections')
        for c in [c for c in self.class_ids.values() if c.section is None]:
            print('    {:>4}: {}'.format(c.cid, c.name))

        # Work with what we have
        self.class_ids = {cid: c for cid, c in self.class_ids.items()
                          if c.section is not None}
        print('')
        print('Parsing deeper for managed Entities with Sections')
        for c in self.class_ids.values():
            print('    {:>9}:  {:>4}: {} -> {}'.format(c.section.section_number,
                                                       c.cid,
                                                       c.name,
                                                       camelcase(c.name)))
            c.deep_parse(self.paragraphs)


if __name__ == '__main__':
    Main().start()
