from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PluginCapability(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    NONE: _ClassVar[PluginCapability]
    VIRTUAL_FILESYSTEM: _ClassVar[PluginCapability]
    ANALYZER: _ClassVar[PluginCapability]
    TRANSFORMER: _ClassVar[PluginCapability]

class DataType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    NULL: _ClassVar[DataType]
    END: _ClassVar[DataType]
    STR: _ClassVar[DataType]
    INT: _ClassVar[DataType]
    DBL: _ClassVar[DataType]
    BIN: _ClassVar[DataType]
    BOOL: _ClassVar[DataType]

class MessageType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    DYNAMIC: _ClassVar[MessageType]
    LOAD: _ClassVar[MessageType]
    ANALYZE: _ClassVar[MessageType]
    FILE_LS: _ClassVar[MessageType]
    FILE_READ: _ClassVar[MessageType]
    SCORE: _ClassVar[MessageType]
    TRANSFORM: _ClassVar[MessageType]
    ERROR: _ClassVar[MessageType]
    TERMINATE: _ClassVar[MessageType]

class EntryType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    UNKNOWN: _ClassVar[EntryType]
    FILE: _ClassVar[EntryType]
    DIRECTORY: _ClassVar[EntryType]
    COMPRESSED: _ClassVar[EntryType]
    NO_INDEX: _ClassVar[EntryType]

class TransformEntryType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
    NOINDEX: _ClassVar[TransformEntryType]
    STRING: _ClassVar[TransformEntryType]
    GROUP: _ClassVar[TransformEntryType]
NONE: PluginCapability
VIRTUAL_FILESYSTEM: PluginCapability
ANALYZER: PluginCapability
TRANSFORMER: PluginCapability
NULL: DataType
END: DataType
STR: DataType
INT: DataType
DBL: DataType
BIN: DataType
BOOL: DataType
DYNAMIC: MessageType
LOAD: MessageType
ANALYZE: MessageType
FILE_LS: MessageType
FILE_READ: MessageType
SCORE: MessageType
TRANSFORM: MessageType
ERROR: MessageType
TERMINATE: MessageType
UNKNOWN: EntryType
FILE: EntryType
DIRECTORY: EntryType
COMPRESSED: EntryType
NO_INDEX: EntryType
NOINDEX: TransformEntryType
STRING: TransformEntryType
GROUP: TransformEntryType

class LoadRequest(_message.Message):
    __slots__ = ["launch_params"]
    LAUNCH_PARAMS_FIELD_NUMBER: _ClassVar[int]
    launch_params: str
    def __init__(self, launch_params: _Optional[str] = ...) -> None: ...

class LoadResponse(_message.Message):
    __slots__ = ["status", "capabilities", "setState", "handlers", "shouldNegotiate"]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CAPABILITIES_FIELD_NUMBER: _ClassVar[int]
    SETSTATE_FIELD_NUMBER: _ClassVar[int]
    HANDLERS_FIELD_NUMBER: _ClassVar[int]
    SHOULDNEGOTIATE_FIELD_NUMBER: _ClassVar[int]
    status: int
    capabilities: _containers.RepeatedScalarFieldContainer[PluginCapability]
    setState: str
    handlers: _containers.RepeatedScalarFieldContainer[str]
    shouldNegotiate: bool
    def __init__(self, status: _Optional[int] = ..., capabilities: _Optional[_Iterable[_Union[PluginCapability, str]]] = ..., setState: _Optional[str] = ..., handlers: _Optional[_Iterable[str]] = ..., shouldNegotiate: bool = ...) -> None: ...

class NegotiateRequest(_message.Message):
    __slots__ = ["seq", "type", "input"]
    SEQ_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    INPUT_FIELD_NUMBER: _ClassVar[int]
    seq: int
    type: DataType
    input: str
    def __init__(self, seq: _Optional[int] = ..., type: _Optional[_Union[DataType, str]] = ..., input: _Optional[str] = ...) -> None: ...

class NegotiateResponse(_message.Message):
    __slots__ = ["seq", "message", "html", "type"]
    SEQ_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    HTML_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    seq: int
    message: str
    html: str
    type: DataType
    def __init__(self, seq: _Optional[int] = ..., message: _Optional[str] = ..., html: _Optional[str] = ..., type: _Optional[_Union[DataType, str]] = ...) -> None: ...

class ReadFileRequest(_message.Message):
    __slots__ = ["path", "name", "mime"]
    PATH_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    MIME_FIELD_NUMBER: _ClassVar[int]
    path: str
    name: str
    mime: str
    def __init__(self, path: _Optional[str] = ..., name: _Optional[str] = ..., mime: _Optional[str] = ...) -> None: ...

class File(_message.Message):
    __slots__ = ["url", "name", "provider_id", "mime_type", "size", "sha256", "md5"]
    URL_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    PROVIDER_ID_FIELD_NUMBER: _ClassVar[int]
    MIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    SIZE_FIELD_NUMBER: _ClassVar[int]
    SHA256_FIELD_NUMBER: _ClassVar[int]
    MD5_FIELD_NUMBER: _ClassVar[int]
    url: str
    name: str
    provider_id: str
    mime_type: str
    size: int
    sha256: str
    md5: str
    def __init__(self, url: _Optional[str] = ..., name: _Optional[str] = ..., provider_id: _Optional[str] = ..., mime_type: _Optional[str] = ..., size: _Optional[int] = ..., sha256: _Optional[str] = ..., md5: _Optional[str] = ...) -> None: ...

class DirectoryEntry(_message.Message):
    __slots__ = ["path", "entry_type", "name", "mime", "provider", "final"]
    PATH_FIELD_NUMBER: _ClassVar[int]
    ENTRY_TYPE_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    MIME_FIELD_NUMBER: _ClassVar[int]
    PROVIDER_FIELD_NUMBER: _ClassVar[int]
    FINAL_FIELD_NUMBER: _ClassVar[int]
    path: str
    entry_type: EntryType
    name: str
    mime: str
    provider: str
    final: bool
    def __init__(self, path: _Optional[str] = ..., entry_type: _Optional[_Union[EntryType, str]] = ..., name: _Optional[str] = ..., mime: _Optional[str] = ..., provider: _Optional[str] = ..., final: bool = ...) -> None: ...

class ListChildrenRequest(_message.Message):
    __slots__ = ["path", "requestSize"]
    PATH_FIELD_NUMBER: _ClassVar[int]
    REQUESTSIZE_FIELD_NUMBER: _ClassVar[int]
    path: str
    requestSize: int
    def __init__(self, path: _Optional[str] = ..., requestSize: _Optional[int] = ...) -> None: ...

class ListChildrenResponse(_message.Message):
    __slots__ = ["entry", "incomplete"]
    ENTRY_FIELD_NUMBER: _ClassVar[int]
    INCOMPLETE_FIELD_NUMBER: _ClassVar[int]
    entry: _containers.RepeatedCompositeFieldContainer[DirectoryEntry]
    incomplete: bool
    def __init__(self, entry: _Optional[_Iterable[_Union[DirectoryEntry, _Mapping]]] = ..., incomplete: bool = ...) -> None: ...

class AnalyzeRequest(_message.Message):
    __slots__ = ["file", "Url"]
    FILE_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    file: File
    Url: str
    def __init__(self, file: _Optional[_Union[File, _Mapping]] = ..., Url: _Optional[str] = ...) -> None: ...

class AnalyzeFinding(_message.Message):
    __slots__ = ["score", "location", "contents", "description"]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    score: int
    location: str
    contents: str
    description: str
    def __init__(self, score: _Optional[int] = ..., location: _Optional[str] = ..., contents: _Optional[str] = ..., description: _Optional[str] = ...) -> None: ...

class AnalyzeResponse(_message.Message):
    __slots__ = ["score", "file", "findings"]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    FILE_FIELD_NUMBER: _ClassVar[int]
    FINDINGS_FIELD_NUMBER: _ClassVar[int]
    score: int
    file: File
    findings: _containers.RepeatedCompositeFieldContainer[AnalyzeFinding]
    def __init__(self, score: _Optional[int] = ..., file: _Optional[_Union[File, _Mapping]] = ..., findings: _Optional[_Iterable[_Union[AnalyzeFinding, _Mapping]]] = ...) -> None: ...

class TransformRequest(_message.Message):
    __slots__ = ["file"]
    FILE_FIELD_NUMBER: _ClassVar[int]
    file: File
    def __init__(self, file: _Optional[_Union[File, _Mapping]] = ...) -> None: ...

class TransformResponse(_message.Message):
    __slots__ = ["file", "url"]
    FILE_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    file: File
    url: str
    def __init__(self, file: _Optional[_Union[File, _Mapping]] = ..., url: _Optional[str] = ...) -> None: ...

class TransformEntry(_message.Message):
    __slots__ = ["type", "uid", "contents", "children", "correlation"]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    UID_FIELD_NUMBER: _ClassVar[int]
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    CHILDREN_FIELD_NUMBER: _ClassVar[int]
    CORRELATION_FIELD_NUMBER: _ClassVar[int]
    type: TransformEntryType
    uid: int
    contents: str
    children: _containers.RepeatedCompositeFieldContainer[TransformEntry]
    correlation: float
    def __init__(self, type: _Optional[_Union[TransformEntryType, str]] = ..., uid: _Optional[int] = ..., contents: _Optional[str] = ..., children: _Optional[_Iterable[_Union[TransformEntry, _Mapping]]] = ..., correlation: _Optional[float] = ...) -> None: ...

class Request(_message.Message):
    __slots__ = ["id", "type", "pluginState", "load", "ls", "cat", "analyze", "transform", "terminate"]
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PLUGINSTATE_FIELD_NUMBER: _ClassVar[int]
    LOAD_FIELD_NUMBER: _ClassVar[int]
    LS_FIELD_NUMBER: _ClassVar[int]
    CAT_FIELD_NUMBER: _ClassVar[int]
    ANALYZE_FIELD_NUMBER: _ClassVar[int]
    TRANSFORM_FIELD_NUMBER: _ClassVar[int]
    TERMINATE_FIELD_NUMBER: _ClassVar[int]
    id: int
    type: MessageType
    pluginState: str
    load: LoadRequest
    ls: ListChildrenRequest
    cat: ReadFileRequest
    analyze: AnalyzeRequest
    transform: TransformRequest
    terminate: TerminateRequest
    def __init__(self, id: _Optional[int] = ..., type: _Optional[_Union[MessageType, str]] = ..., pluginState: _Optional[str] = ..., load: _Optional[_Union[LoadRequest, _Mapping]] = ..., ls: _Optional[_Union[ListChildrenRequest, _Mapping]] = ..., cat: _Optional[_Union[ReadFileRequest, _Mapping]] = ..., analyze: _Optional[_Union[AnalyzeRequest, _Mapping]] = ..., transform: _Optional[_Union[TransformRequest, _Mapping]] = ..., terminate: _Optional[_Union[TerminateRequest, _Mapping]] = ...) -> None: ...

class TerminateRequest(_message.Message):
    __slots__ = ["id"]
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class TerminateResponse(_message.Message):
    __slots__ = ["id"]
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class Error(_message.Message):
    __slots__ = ["trace", "message"]
    TRACE_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    trace: str
    message: str
    def __init__(self, trace: _Optional[str] = ..., message: _Optional[str] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ["id", "type", "pluginState", "incomplete", "err", "load", "ls", "cat", "analyze", "transform", "terminate"]
    ID_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PLUGINSTATE_FIELD_NUMBER: _ClassVar[int]
    INCOMPLETE_FIELD_NUMBER: _ClassVar[int]
    ERR_FIELD_NUMBER: _ClassVar[int]
    LOAD_FIELD_NUMBER: _ClassVar[int]
    LS_FIELD_NUMBER: _ClassVar[int]
    CAT_FIELD_NUMBER: _ClassVar[int]
    ANALYZE_FIELD_NUMBER: _ClassVar[int]
    TRANSFORM_FIELD_NUMBER: _ClassVar[int]
    TERMINATE_FIELD_NUMBER: _ClassVar[int]
    id: int
    type: MessageType
    pluginState: str
    incomplete: bool
    err: Error
    load: LoadResponse
    ls: ListChildrenResponse
    cat: File
    analyze: AnalyzeResponse
    transform: TransformResponse
    terminate: TerminateResponse
    def __init__(self, id: _Optional[int] = ..., type: _Optional[_Union[MessageType, str]] = ..., pluginState: _Optional[str] = ..., incomplete: bool = ..., err: _Optional[_Union[Error, _Mapping]] = ..., load: _Optional[_Union[LoadResponse, _Mapping]] = ..., ls: _Optional[_Union[ListChildrenResponse, _Mapping]] = ..., cat: _Optional[_Union[File, _Mapping]] = ..., analyze: _Optional[_Union[AnalyzeResponse, _Mapping]] = ..., transform: _Optional[_Union[TransformResponse, _Mapping]] = ..., terminate: _Optional[_Union[TerminateResponse, _Mapping]] = ...) -> None: ...
