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
import re
import json
from text import ascii_only


class TableList(object):
    def __init__(self, tables=None):
        self._docx_tables = tables
        self._tables = list()

    def __getitem__(self, item):
        return self._tables[item]  # delegate to li.__getitem__

    def __iter__(self):
        for table in self._tables:
            yield table

    def __len__(self):
        return len(self._tables)

    def add(self, table):
        assert isinstance(table, Table)
        self._tables.append(table)
        return self

    @staticmethod
    def create(tables):
        tbl_list = TableList(tables=tables)

        # Decode is about 2-1/2 minutes
        for idx, tbl in enumerate(tables):
            table = Table.create(idx, tbl)
            tbl_list.add(table)

        return tbl_list

    def get(self, index):
        return self._tables[index]

    def size(self):
        return len(self._tables)

    # def save(self, filepath):
    #     data = json.dumps(self.as_dict_list(), indent=2, separators=(',', ': '))
    #     with open(filepath, 'w') as f:
    #         f.write(data)
    #
    # def load(self, filepath):
    #     self._tables = list()
    #     with open(filepath, 'r') as f:
    #         data = json.load(f)
    #         for head in data:
    #             # section = SectionHeading()
    #             # section.paragraph_number = head['paragraph_number']
    #             # section.style_name = head['style_name']
    #             # section.section_number = head['section_number']
    #             # section.title = head['title']
    #             # section.section_points = head['section_points']
    #
    #             # self._tables.append(table)
    #             pass
    #
    # def as_dict_list(self):
    #     return [
    #         {
    #             'paragraph_number': head.paragraph_number,
    #             'style_name': head.style_name,
    #             'section_number': head.section_number,
    #             'title': head.title,
    #             'section_points': head.section_points,
    #         } for head in self._tables]


class Table(object):
    """
    Minimal table info

    There are several formats tables come in:

    Table with a numbered title. Detect by row-0 having the word  'Table '
    or 'Table-' as one or more cell text:

        Table x-y:  Title                   <- Table Title with number following
        Heading 1  | Heading 2  | ...       <- Column Headings
        Row-Data-1 | Row-Data-2 | ...
        ...

    Table with column headings and no title. Look into first row of 'cells'
    and examine the cell[n].paragraphs[n].style and see if it its style starts
    with 'Table_head' or 'Attribute follower':

        Heading 1  | Heading 2  | ...       <- Column Headings (only)
        Row-Data-1 | Row-Data-2 | ...
        ...

    Table with without heading.  Row has 'Table_text' or 'Attribute style
        Row-Data-1 | Row-Data-2 | ...
        ...
    """
    def __init__(self):
        self.rows = []
        self.heading = None
        self.doc_table_number = None        # Index into document table list
        self.num_columns = 0
        self.full_title = None
        self.table_number = None            # From full title string
        self.short_title = None             # From full title string

    def __str__(self):
        return 'Table {}: columns: {}, rows: {}, {}'.format(self.doc_table_number,
                                                            self.num_columns,
                                                            len(self.rows),
                                                            self.full_title)

    @staticmethod
    def create(index, doc_table):
        table = Table()
        table.doc_table_number = index

        try:
            table.num_columns = len(doc_table.columns)

            for row_num, row in enumerate(doc_table.rows):
                text = (ascii_only(cell.text) for cell in row.cells)

                # Establish the mapping based on the first row
                # headers; these will become the keys of our dictionary

                if table.heading is None:
                    # Title or just a heading?
                    text_tuple = tuple(text)
                    title = next((t for t in text_tuple if 'Table ' in t or 'Table-' in t),None)

                    if title is not None:
                        # Table with a numbered title
                        if 'Table-' in title:
                            title = 'Table ' + title[6:]

                        table.full_title = title.strip()
                        tparts = re.findall(r"[^\s']+", table.full_title)
                        table.table_number = tparts[1]
                        table.short_title = ' '.join(tparts[2:])
                        continue

                    elif any('table_head' in c.paragraphs[0].style.name.lower() for c in row.cells) \
                            or any('attribute follower' in c.paragraphs[0].style.name.lower() for c in row.cells):
                        # Table with column headings and no title
                        table.heading = text_tuple
                        continue

                    elif all(text_tuple[0].strip().lower() == t.strip().lower() for t in text_tuple[1:]):
                        table.full_title = next((t.strip() for t in text_tuple if len(t.strip())), text_tuple[0])
                        table.short_title = table.full_title
                        continue

                    elif any('table_text' in c.paragraphs[0].style.name.lower() for c in row.cells) \
                            or any('attribute' in c.paragraphs[0].style.name.lower() for c in row.cells):
                        # Table with without heading
                        table.heading = tuple('col-{}'.format(n) for n in range(1, table.num_columns+1))
                        continue

                    else:
                        # Default type (low importance usually)
                        table.heading = text_tuple
                        continue

                # elif row_num == 1 and table.full_title is not None:
                #     text_tuple = tuple(text)
                #     table.heading = text_tuple
                #     continue

                # Construct a dictionary for this row, mapping
                # keys to values for this row
                # row_data = dict(zip(table.heading, ascii_only(text)))
                row_data = dict(zip(table.heading, text))

                table.rows.append(row_data)

            return table

        except Exception as _e:
            raise
