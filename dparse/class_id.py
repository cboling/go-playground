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
import re
from text import ascii_only
from section import SectionHeading
from tables import Table
from transitions import Machine
from contents import *


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

    def keys(self):
        return self._entities.keys()

    def values(self):
        return self._entities.values()

    def items(self):
        return self._entities.items()

    def has_key(self, k):
        return k in self._entities

    def add(self, class_id):
        assert isinstance(class_id, ClassId), 'Invalid type'
        assert class_id.cid not in self._entities, \
            'Entity already exists: {}'.format(class_id)
        self._entities[class_id.cid] = class_id
        return self

    @staticmethod
    def parse_sections(sections, cid_section):
        """
        Look for a specific section and determine the table it is in.

        After finding the section, this is more focused on the latest
        G.988 document where the list is in the one and only table
        """
        cid_list = ClassIdList()
        cid_heading_section = sections.find_section(cid_section)

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

                        cid.section = sections.find_section_by_name(cid.name)

                    except ValueError as _e:
                        pass        # Expected for reserved range statements

                    except Exception as _e:
                        raise              # Not expected

        return cid_list


class ClassId(object):
    """ Managed Entity Class Information """
    STATES = ['initial', 'description', 'relationships', 'attributes',
              'operations', 'optionals', 'notificatons', 'alarms', 'avcs',
               'tests', 'complete', 'failure']

    TRANSITIONS = [
        # While in initial 'basic' state
        {'trigger': 'normal', 'source': 'basic', 'dest': 'description'},
        {'trigger': 'relationship', 'source': 'initial', 'dest': 'relationships'},
        {'trigger': 'attribute', 'source': 'initial', 'dest': 'attributes'},

        # While in 'description' state
        {'trigger': 'normal', 'source': 'description', 'dest': 'description'},
        {'trigger': 'relationship', 'source': 'description', 'dest': 'relationships'},
        {'trigger': 'attribute', 'source': 'description', 'dest': 'attribute'},
        # {'trigger': '', 'source': 'description', 'dest': ''},

        # While in 'relationships' state
        {'trigger': 'normal', 'source': 'relationships', 'dest': 'relationships'},
        {'trigger': 'attribute', 'source': 'relationships', 'dest': 'attribute'},
        # {'trigger': '', 'source': 'relationships', 'dest': ''},
        # {'trigger': '', 'source': 'relationships', 'dest': ''},

        # While in 'attributes' state
        {'trigger': 'normal', 'source': 'attributes', 'dest': 'attributes'},
        {'trigger': 'operation', 'source': 'attributes', 'dest': 'operations'},
        # {'trigger': '', 'source': 'attributes', 'dest': ''},
        # {'trigger': '', 'source': 'attributes', 'dest': ''},

        # While in 'operations' state
        {'trigger': 'normal', 'source': 'operations', 'dest': 'operations'},
        {'trigger': 'optional', 'source': 'operations', 'dest': 'optionals'},
        {'trigger': 'notification', 'source': 'operations', 'dest': 'notifications'},
        # {'trigger': '', 'source': 'operations', 'dest': ''},
        # {'trigger': '', 'source': 'operations', 'dest': ''},

        # While in 'optionals' state
        {'trigger': 'normal', 'source': 'optionals', 'dest': 'optionals'},
        {'trigger': 'notification', 'source': 'optionals', 'dest': 'notifications'},
        # {'trigger': '', 'source': 'optionals', 'dest': ''},

        # While in 'notifications' state
        {'trigger': 'normal', 'source': 'notifications', 'dest': 'notifications'},
        {'trigger': 'alarm', 'source': 'notifications', 'dest': 'alarms'},
        {'trigger': 'avc', 'source': 'notifications', 'dest': 'avcs'},
        {'trigger': 'test', 'source': 'notifications', 'dest': 'tests'},

        # While in 'alarms' state
        {'trigger': 'normal', 'source': 'alarms', 'dest': 'alarms'},
        {'trigger': 'avc', 'source': 'alarms', 'dest': 'avcs'},
        {'trigger': 'test', 'source': 'alarms', 'dest': 'tests'},
        # {'trigger': '', 'source': 'alarms', 'dest': ''},

        # While in 'avc' state
        {'trigger': 'normal', 'source': 'avcs', 'dest': 'source'},
        {'trigger': 'alarm', 'source': 'avcs', 'dest': 'alarms'},
        {'trigger': 'test', 'source': 'avcs', 'dest': 'tests'},
        # {'trigger': '', 'source': 'avcs', 'dest': ''},

        # While in 'tests' state
        {'trigger': 'normal', 'source': 'tests', 'dest': 'tests'},
        {'trigger': 'alarm', 'source': 'tests', 'dest': 'alarms'},
        {'trigger': 'avc', 'source': 'tests', 'dest': 'avcs'},
        # {'trigger': '', 'source': 'tests', 'dest': ''},

        # Do wildcard 'complete' trigger last so it covers all previous states
        {'trigger': 'complete', 'source': '*', 'dest': 'complete'},
        {'trigger': 'failure', 'source': '*', 'dest': 'failure'},
    ]

    def __init__(self):
        self.cid = None                   # Class Id
        self.name = None                  # Title
        self.section = None               # Document section

        self.parser = initial_parser
        self.machine = Machine(model=self, states=ClassId.STATES,
                               transitions=ClassId.TRANSITIONS,
                               initial='initial',
                               queued=True,
                               name='{}-{}'.format(self.name, self.cid))

        # Following hold lists of paragraph numbers
        self.description = list()         # Description
        self.relationships = list()       # Relationships paragraph (if any)

        # Following hold lists of associated objects
        self.attributes = list()          # Ordered list of attributes
        self.operations = set()           # Mandatory operations/message-types
        self.optional_operations = set()  # Allowed operations/message-types
        self.alarms = None                # Alarm list (if any)
        self.avcs = None                  # Attribute Value Change list (if any)
        self.test_results = None          # Test Results (if any)
        self.hidden = False               # Not reported or ignore in MIB upload

    def __str__(self):
        return 'ClassId: {}: {}'.format(self.cid, self.name)

    def deep_parse(self, paragraphs):
        """ Fill out detailed class information """
        if self.section is None:
            return self

        for content in self.section.contents:
            try:
                if isinstance(content, int):
                    # Paragraph number
                    trigger = self.parser(content, paragraphs)

                elif isinstance(content, Table):
                    # Table object
                    # TODO: trigger = self.parser(content, paragraphs)
                    pass

                else:
                    raise NotImplementedError('Unknown content type: {}'.
                                              format(type(content)))
            except Exception as e:
                self.failure()
                print("FAILURE: Exited deep parsing: '{}'".format(e.message))
                break

        return self

    def on_enter_initial(self):
        self.parser = initial_parser
        pass

    def on_enter_description(self):
        self.parser = description_parser
        pass

    def on_enter_relationships(self):
        self.parser = relationships_parser
        pass

    def on_enter_attributes(self):
        self.parser = attributes_parser
        pass

    def on_enter_operations(self):
        self.parser = operations_parser
        pass

    def on_enter_optionals(self):
        self.parser = optionals_parser
        pass

    def on_enter_notifications(self):
        self.parser = notifications_parser
        pass

    def on_enter_alarms(self):
        self.parser = alarms_parser
        pass

    def on_enter_avcs(self):
        self.parser = avcs_parser
        pass

    def on_enter_tests(self):
        self.parser = tests_parser
        pass

    def on_enter_complete(self):
        self.parser = None
        pass

    def on_enter_failure(self):
        self.parser = None
        pass
