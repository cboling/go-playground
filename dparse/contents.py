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
from text import *


def initial_parser(content, paragraphs):
    """
    Parse content for the 'initial' state of an ME.

    Possible transitions include:
        - Relationships  - if 'Relationships' paragraph header found
        - Attributes     - if 'Attributes' paragraph header found
        - Description    - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_relationships_header(paragraph):
            return 'relationship', None

        elif is_attributes_header(paragraph):
            return 'attribute', None

        elif is_description_style(paragraph.style):
            return 'description', ascii_only(paragraph.text)

        elif is_normal_style(paragraph.style):
            return 'normal', ascii_only(paragraph.text)

    return 'failure', None


def description_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    Possible transitions include:
        - Relationships  - if 'Relationships' paragraph header found
        - Attributes     - if 'Attributes' paragraph header found
        - Description    - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_relationships_header(paragraph):
            return 'relationship', None

        elif is_attributes_header(paragraph):
            return 'attribute', None

        elif is_description_text(paragraph):
            return 'normal', ascii_only(paragraph.text)
    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def relationships_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    Possible transitions include:
        - Attributes     - if 'Attributes' paragraph header found
        - Relationships  - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_attributes_header(paragraph):
            return 'attribute', None

        elif is_relationships_text(paragraph):
            return 'normal', ascii_only(paragraph.text)
    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def attributes_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    Possible transitions include:
        - Operations    - if 'Operations' paragraph header found
        - Attributes    - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_operations_header(paragraph):
            return 'operation', None

        elif is_attribute_text(paragraph):
            return 'normal', ascii_only(paragraph.text)
    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def operations_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    Possible transitions include:
        - Notifications - if 'Notifications paragraph header found
        - Operation     - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_notifications_header(paragraph):
            return 'notification', None

        elif is_operations_text(paragraph):
            return 'normal', ascii_only(paragraph.text)
    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def notifications_parser(content, paragraphs):
    """
    Parse content for this state of an ME.  Unlike other
    ME sections, this will have either the phrase 'None'
    or it may have alarm, avc, or test results subsections.

    Possible transitions include:
        - Alarms        - if 'Alarms paragraph header found
        - AVCs          - if 'Attribute Value Changes paragraph header found
        - Tests         - if 'Test Results paragraph header found
        - Notifications - if 'normal' text style

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_avcs_header(paragraph):
            return 'avc', None

        elif is_alarms_header(paragraph):
            return 'alarm', ascii_only(paragraph.text)

        elif is_tests_header(paragraph):
            return 'test', ascii_only(paragraph.text)

        elif is_notifications_text(paragraph):
            return 'normal', ascii_only(paragraph.text)

    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def alarms_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_avcs_header(paragraph):
            return 'avc', None

        elif is_tests_header(paragraph):
            return 'test', ascii_only(paragraph.text)

        elif is_alarms_text(paragraph):
            return 'normal', ascii_only(paragraph.text)

    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'


def avcs_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_alarms_header(paragraph):
            return 'alarm', ascii_only(paragraph.text)

        elif is_tests_header(paragraph):
            return 'test', ascii_only(paragraph.text)

        elif is_avcs_text(paragraph):
            return 'normal', ascii_only(paragraph.text)

    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')
    return 'failure'


def tests_parser(content, paragraphs):
    """
    Parse content for this state of an ME

    :param content: (int) or (Table)
    :param paragraphs: (list) Docx Paragraphs

    :return: (str, str) Next state, Associated text (if any)
    """
    if isinstance(content, int):
        paragraph = paragraphs[content]

        if is_avcs_header(paragraph):
            return 'avc', None

        elif is_alarms_header(paragraph):
            return 'alarm', ascii_only(paragraph.text)

        elif is_tests_text(paragraph):
            return 'normal', ascii_only(paragraph.text)

    else:
        # TODO: Implement if needed, otherwise remove and fall through
        raise NotImplementedError('Table support')

    return 'failure'

