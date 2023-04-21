import json
import torch
import re

print("Torch version " + torch.__version__)

from langchain.text_splitter import RecursiveCharacterTextSplitter
from transformers import pipeline, AutoTokenizer

from langchain import PromptTemplate, LLMChain, HuggingFacePipeline
from langchain.docstore.document import Document

ACTION_ITEM_AI_MODEL = "declare-lab/flan-alpaca-gpt4-xl"
TOKEN_LIMIT = 480

TARGET_DEVICE = torch.device('cuda' if torch.cuda.is_available() else ("mps" if torch.backends.mps.is_available() else "cpu"))
print(f"Using device: {TARGET_DEVICE}")
ACTION_POINT_PIPELINE = pipeline(model=ACTION_ITEM_AI_MODEL, tokenizer=ACTION_ITEM_AI_MODEL, device=TARGET_DEVICE,
                                 max_length=200)


def action_items(transcript):
    print(f"inside summarise:\n{json.dumps(transcript)}")
    tokenizer = AutoTokenizer.from_pretrained(ACTION_ITEM_AI_MODEL)
    transcript_documents = split_text_into_documents(tokenizer, [x['text'] for x in transcript])

    print("Transcript split into " + str(len(transcript_documents)) + " documents")
    prompt = PromptTemplate(
        input_variables=["context"],
        template="List all tasks that a person promised to do from the following input, fully explain each task. If do not list tasks if they are unclear. If there are no action points, say \"None\".\n\nInput:\n{context}")
    merge_prompt = PromptTemplate(
        input_variables=["context"],
        template="The following input has a list of tasks. Make a new list with the most important tasks, filter out similar tasks.\n\nInput:\n{context}")

    llm = HuggingFacePipeline(pipeline=ACTION_POINT_PIPELINE)


    transcript_action_items = []
    for transcript_document in transcript_documents:
        llm_context_chain = LLMChain(llm=llm, prompt=prompt)
        output = llm_context_chain.run(instruction="", context=transcript_document.page_content)

        transcript_action_items += re.split(r"\d+\. ", output)
        print(f"Primary output run: {output}")

    transcript_action_items = [x.strip() for x in transcript_action_items if x != "None" and x != ""]

    print(f"Transcript Action Points: {json.dumps(transcript_action_items)}")
    action_chunks = split_text_into_documents(tokenizer, transcript_action_items)

    run_count = 0
    while len(action_chunks) > 1:
        new_action_list = []
        for action_chunk in action_chunks:
            llm_context_chain = LLMChain(llm=llm, prompt=merge_prompt)
            output = llm_context_chain.run(context=action_chunk.page_content)
            print(f"Intermediate output run {run_count}: {output}")
            new_action_list.append(output)
        action_chunks = split_text_into_documents(tokenizer, new_action_list)
        run_count += 1

    print(f"Merged output: {action_chunks[0].page_content}")
    action_list = re.split(r"(\d+\. |\n)", action_chunks[0].page_content)
    action_list = [x.strip() for x in action_list]
    action_list = [x for x in action_list if x != "" and re.search(r"^\d+\.", x) is None]
    print(f"Action Points: {json.dumps(action_list[1:])}")

    return action_list[1:]


def split_text_into_documents(tokenizer, transcript):
    dialogue = ""
    transcript_documents = []
    text_splitter = RecursiveCharacterTextSplitter.from_huggingface_tokenizer(tokenizer, separators=["\n", ".", " "],
                                                                              chunk_size=TOKEN_LIMIT, chunk_overlap=0)
    for line in transcript:
        dialogue = dialogue + line + "\n"
    chunks = text_splitter.split_text(dialogue)
    for chunk in chunks:
        tokens = len(tokenizer.encode(chunk))
        print(f"Chunk {tokens}:\n{chunk}\n\n")
        transcript_documents.append(Document(page_content=chunk))
    return transcript_documents
