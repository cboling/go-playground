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

from __future__ import (
    absolute_import, division, print_function, unicode_literals
)

from docx import Document

from section import SectionHeading, SectionList
from class_id import ClassIdList, ClassId
from tables import TableList, Table
from text import ascii_only


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

    def parse_me_class_ids(self):
        """
        Look for a specific section and determine the table it is in.

        After finding the section, this is more focused on the latest
        G.988 document where the list is in the one and only table
        """
        cid_list = ClassIdList()
        cid_heading_section = self.find_section(self.args['me-class-section'])

        if cid_heading_section is not None:
            cid_table = next((c for c in cid_heading_section.contents
                              if isinstance(c, Table)), None)
            if cid_table is not None:
                headings = cid_table.heading
                for row in cid_table.rows:
                    try:
                        cid = ClassId()
                        cid.cid = int(row.get(headings[0]))
                        cid.name = row.get(headings[1])
                        cid_list.add(cid)

                        cid.section = self.find_section_by_name(cid.name)

                    except ValueError as _e:
                        pass        # Expected for reserved range statements

                    except Exception as _e:
                        pass              # Not expected

        return cid_list

    def find_section(self, section_number):
        entry = next((s for s in self.sections if s.section_number == section_number), None)

        if entry is not None:
            return entry

        raise KeyError('Section {} not found'.format(section_number))

    def find_section_by_name(self, name):
        name_lower = name.replace(' ', '').lower()
        return next((s for s in self.sections
                    if s.title.replace(' ', '').lower() == name_lower), None)

    def get_tables(self, tables):
        return TableList.create(tables)

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

        self.class_ids = self.parse_me_class_ids()
        print('Found {} ME Class ID entries'.format(len(self.class_ids)))


if __name__ == '__main__':
    Main().start()
