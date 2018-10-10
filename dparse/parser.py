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
        """ Look for a specific section and determine the table it is in"""
        cid_list = ClassIdList()
        cid_heading_section = self.find_section(self.args['me-class-section'])

        # The first paragraph number for a section is the heading paragraph
        # heading_paragraph = self.paragraphs[cid_heading_section.paragraph_number[0]]
        # next = cid_heading_section.paragraph_number + 1
        # text = list()

        for pnum in cid_heading_section.paragraph_numbers:
            next_para = self.paragraphs[pnum]

            print('Paragraph : {}'.format(pnum))
            print('  Style   : {}'.format(next_para.style.name))
            print('  Contents: {}'.format(ascii_only(next_para.text)))

            if next_para.style.builtin and 'heading ' in next_para.style.name.lower():
                break

        print('Found {} text sections'.format(len(text)))

        return cid_list

    def find_section(self, section_number):
        entry = next((s for s in self.sections if s.section_number == section_number), None)

        if entry is not None:
            return entry

        raise KeyError('Section {} not found'.format(section_number))

    def get_tables(self, tables):
        return TableList.create(tables)


    def start(self):
        print("Loading ITU Document '{}' and parsed Section Header file '{}'".
              format(self.args['itu'], self.args['sections']))

        document = self.load_itu_document()
        self.sections = SectionList()
        self.sections.load(self.args['sections'])

        self.paragraphs = document.paragraphs
        # doc_sections = document.sections
        # styles = document.styles
        self.body = document.element.body

        tables = self.get_tables(document.tables)

        print('Extracting ME Class ID values')

        self.class_ids = self.parse_me_class_ids()
        print('Found {} ME Class ID entries'.format(len(self.class_ids)))


if __name__ == '__main__':
    Main().start()
